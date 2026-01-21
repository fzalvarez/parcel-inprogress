package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	manifestusecase "ms-parcel-core/internal/parcel/parcel_manifest/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type ManifestPreviewRequest struct {
	VehicleID           string `json:"vehicle_id" binding:"required,uuid"`
	OriginOfficeID      string `json:"origin_office_id" binding:"required,uuid"`
	DestinationOfficeID string `json:"destination_office_id" binding:"required,uuid"`
}

type ManifestHandler struct {
	buildUC *manifestusecase.BuildManifestPreviewUseCase
}

func NewManifestHandler(buildUC *manifestusecase.BuildManifestPreviewUseCase) *ManifestHandler {
	return &ManifestHandler{buildUC: buildUC}
}

// PreviewPost godoc
// @Summary Construir preview de manifiesto (POST)
// @Description Construye un manifiesto virtual (preview) basado en envíos pendientes entre oficinas. Acepta vehículo, oficina de origen y destino. El preview incluye listado de envíos, totales (cantidad, peso, volumen) y detalles de rutas.
// @Tags Manifests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param payload body ManifestPreviewRequest true "Solicitud de preview con IDs de vehículo y oficinas"
// @Success 200 {object} handler.AnyDataEnvelope "Preview de manifiesto generado"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: UUID inválido, payload malformado o parámetros faltantes"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Vehículo u oficina no encontrados"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /manifests/preview [post]
func (h *ManifestHandler) PreviewPost(c *gin.Context) {
	var req ManifestPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}
	// binding ya valida UUID, pero reforzamos por consistencia
	if _, err := uuid.Parse(req.VehicleID); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "vehicle_id inválido", map[string]any{"field": "vehicle_id"}))
		return
	}
	if _, err := uuid.Parse(req.OriginOfficeID); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "origin_office_id inválido", map[string]any{"field": "origin_office_id"}))
		return
	}
	if _, err := uuid.Parse(req.DestinationOfficeID); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "destination_office_id inválido", map[string]any{"field": "destination_office_id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	prev, err := h.buildUC.Execute(c.Request.Context(), manifestusecase.BuildManifestPreviewInput{
		TenantID:            tenant,
		VehicleID:           req.VehicleID,
		OriginOfficeID:      req.OriginOfficeID,
		DestinationOfficeID: req.DestinationOfficeID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": prev})
}

// PreviewGet godoc
// @Summary Construir preview de manifiesto (GET)
// @Description Construye un manifiesto virtual (preview) basado en parámetros de query. Acepta vehículo, oficina de origen y destino. El preview incluye listado de envíos, totales (cantidad, peso, volumen) y detalles de rutas.
// @Tags Manifests
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param vehicle_id query string true "UUID del vehículo" Format(uuid)
// @Param origin_office_id query string true "UUID de la oficina de origen" Format(uuid)
// @Param destination_office_id query string true "UUID de la oficina de destino" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope "Preview de manifiesto generado"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: UUID inválido o parámetros faltantes"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Vehículo u oficina no encontrados"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /manifests/preview [get]
func (h *ManifestHandler) PreviewGet(c *gin.Context) {
	vehicleID := strings.TrimSpace(c.Query("vehicle_id"))
	originOfficeID := strings.TrimSpace(c.Query("origin_office_id"))
	destinationOfficeID := strings.TrimSpace(c.Query("destination_office_id"))

	if vehicleID == "" || originOfficeID == "" || destinationOfficeID == "" {
		_ = c.Error(apperror.NewBadRequest("validation_error", "faltan parámetros", map[string]any{"required": []string{"vehicle_id", "origin_office_id", "destination_office_id"}}))
		return
	}
	if _, err := uuid.Parse(vehicleID); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "vehicle_id inválido", map[string]any{"field": "vehicle_id"}))
		return
	}
	if _, err := uuid.Parse(originOfficeID); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "origin_office_id inválido", map[string]any{"field": "origin_office_id"}))
		return
	}
	if _, err := uuid.Parse(destinationOfficeID); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "destination_office_id inválido", map[string]any{"field": "destination_office_id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	prev, err := h.buildUC.Execute(c.Request.Context(), manifestusecase.BuildManifestPreviewInput{
		TenantID:            tenant,
		VehicleID:           vehicleID,
		OriginOfficeID:      originOfficeID,
		DestinationOfficeID: destinationOfficeID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": prev})
}
