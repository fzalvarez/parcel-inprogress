package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	coreusecase "ms-parcel-core/internal/parcel/parcel_core/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

const parcelSummaryTrackingLimit = 20

type ParcelSummaryHandler struct {
	uc *coreusecase.GetParcelSummaryUseCase
}

func NewParcelSummaryHandler(uc *coreusecase.GetParcelSummaryUseCase) *ParcelSummaryHandler {
	return &ParcelSummaryHandler{uc: uc}
}

// Get godoc
// @Summary Resumen operativo
// @Description Devuelve parcel + items + payment + tracking
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/summary [get]
func (h *ParcelSummaryHandler) Get(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	parcelID, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	out, err := h.uc.Execute(c.Request.Context(), tenant, parcelID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Parcel: reutiliza el mismo shape que GET /parcels/:id devolviendo el domain.Parcel tal como ya lo hace ParcelHandler.
	parcel := out.Parcel

	items := make([]ParcelItemResponse, 0, len(out.Items))
	for _, it := range out.Items {
		items = append(items, ParcelItemResponse{
			ID:          it.ID,
			ParcelID:    it.ParcelID,
			Description: it.Description,
			Quantity:    it.Quantity,
			WeightKg:    it.WeightKg,
			UnitPrice:   it.UnitPrice,
			ContentType: it.ContentType,
			Notes:       it.Notes,
			CreatedAt:   it.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	var payment *ParcelPaymentResponse
	if out.Payment != nil {
		var paidAtStr *string
		if out.Payment.PaidAt != nil {
			s := out.Payment.PaidAt.UTC().Format(time.RFC3339)
			paidAtStr = &s
		}
		payment = &ParcelPaymentResponse{
			ID:           out.Payment.ID,
			ParcelID:     out.Payment.ParcelID,
			PaymentType:  string(out.Payment.PaymentType),
			Currency:     string(out.Payment.Currency),
			Amount:       out.Payment.Amount,
			Notes:        out.Payment.Notes,
			Status:       string(out.Payment.Status),
			CreatedAt:    out.Payment.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:    out.Payment.UpdatedAt.UTC().Format(time.RFC3339),
			PaidAt:       paidAtStr,
			PaidByUserID: out.Payment.PaidByUserID,
		}
	}

	tracking := make([]gin.H, 0, len(out.Tracking))
	for _, e := range out.Tracking {
		tracking = append(tracking, gin.H{
			"id":          e.ID,
			"parcel_id":   e.ParcelID,
			"event_type":  e.EventType,
			"occurred_at": e.OccurredAt.UTC().Format(time.RFC3339),
			"user_id":     e.UserID,
			"user_name":   e.UserName,
			"metadata":    e.Metadata,
		})
	}

	// Asegurar límite 20 también en handler por consistencia
	if len(tracking) > parcelSummaryTrackingLimit {
		tracking = tracking[:parcelSummaryTrackingLimit]
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"parcel":   parcel,
			"items":    items,
			"payment":  payment,
			"tracking": tracking,
			"meta": gin.H{
				"tracking_limit": parcelSummaryTrackingLimit,
			},
		},
	})
}

// Helpers mínimos para reutilizar DTOs existentes sin cambiar shapes.

func mapParcelToResponse(p any) any {
	// ParcelHandler ya expone el shape; aquí devolvemos el mismo objeto domain.Parcel mapeado por el handler existente.
	// Reutilizamos la misma estructura que Create/GetByID están usando.
	parcel := p
	_ = parcel
	return p
}

func mapItemsToResponse(items any) any {
	return items
}

func mapPaymentToResponse(payment any) any {
	return payment
}

func mapTrackingToResponse(events any) any {
	return events
}
