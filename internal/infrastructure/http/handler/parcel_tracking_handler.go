package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_tracking/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type ParcelTrackingHandler struct {
	listUC *usecase.ListTrackingUseCase
}

func NewParcelTrackingHandler(listUC *usecase.ListTrackingUseCase) *ParcelTrackingHandler {
	return &ParcelTrackingHandler{listUC: listUC}
}

// ListByParcelID godoc
// @Summary Listar historial completo de tracking del envío
// @Description Lista todos los eventos y cambios de estado registrados para un envío en orden cronológico. Incluye tipo de evento, timestamp, usuario responsable y metadata adicional. Útil para auditoría, debugging y reconstrucción de la cadena de custodia.
// @Tags ParcelTracking
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope "Historial de tracking del envío"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío no encontrado"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /parcels/{id}/tracking [get]
func (h *ParcelTrackingHandler) ListByParcelID(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(idStr); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	evs, err := h.listUC.Execute(c.Request.Context(), tenant, idStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	out := make([]map[string]any, 0, len(evs))
	for _, ev := range evs {
		out = append(out, map[string]any{
			"id":          ev.ID.String(),
			"parcel_id":   ev.ParcelID,
			"event_type":  ev.EventType,
			"occurred_at": ev.OccurredAt.UTC().Format(time.RFC3339),
			"user_id":     ev.UserID,
			"user_name":   ev.UserName,
			"metadata":    ev.Metadata,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
}
