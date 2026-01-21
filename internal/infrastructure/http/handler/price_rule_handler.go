package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	pricingdomain "ms-parcel-core/internal/parcel/parcel_pricing/domain"
	pricingusecase "ms-parcel-core/internal/parcel/parcel_pricing/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type PriceRuleRequest struct {
	ShipmentType        string  `json:"shipment_type" binding:"required"`
	OriginOfficeID      string  `json:"origin_office_id" binding:"required"`
	DestinationOfficeID string  `json:"destination_office_id" binding:"required"`
	Unit                string  `json:"unit" binding:"required,oneof=PER_KG PER_ITEM"`
	Price               float64 `json:"price" binding:"required"`
	Currency            string  `json:"currency" binding:"required,oneof=PEN USD"`
	Priority            int     `json:"priority" binding:"omitempty,min=0,max=100"`
	Active              bool    `json:"active"`
}

type PriceRuleResponse struct {
	ID                  string  `json:"id"`
	ShipmentType        string  `json:"shipment_type"`
	OriginOfficeID      string  `json:"origin_office_id"`
	DestinationOfficeID string  `json:"destination_office_id"`
	Unit                string  `json:"unit"`
	Price               float64 `json:"price"`
	Currency            string  `json:"currency"`
	Priority            int     `json:"priority"`
	Active              bool    `json:"active"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}

type PriceRuleHandler struct {
	createUC *pricingusecase.CreatePriceRuleUseCase
	updateUC *pricingusecase.UpdatePriceRuleUseCase
	listUC   *pricingusecase.ListPriceRulesUseCase
}

func NewPriceRuleHandler(createUC *pricingusecase.CreatePriceRuleUseCase, updateUC *pricingusecase.UpdatePriceRuleUseCase, listUC *pricingusecase.ListPriceRulesUseCase) *PriceRuleHandler {
	return &PriceRuleHandler{createUC: createUC, updateUC: updateUC, listUC: listUC}
}

// Create godoc
// @Summary Crear regla de precios
// @Description Crea una nueva regla de precios para el tenant. Soporta comodines (*) en ShipmentType, OriginOfficeID y DestinationOfficeID. La prioridad (0-100) determina el orden de evaluación en búsquedas jerárquicas: específicas primero, luego comodines. Ejemplos: "STANDARD", "*" para cualquier tipo; "12345" (UUID), "*" para cualquier oficina.
// @Tags Pricing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param payload body PriceRuleRequest true "Solicitud de creación de regla con campos de envío, precio y prioridad"
// @Success 200 {object} handler.AnyDataEnvelope "Regla de precios creada exitosamente"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: payload malformado, valores inválidos o formato incorrecto"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 409 {object} handler.ErrorResponse "Conflicto: regla duplicada o combinación de parámetros duplicada"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /pricing/rules [post]
func (h *PriceRuleHandler) Create(c *gin.Context) {
	var req PriceRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	out, err := h.createUC.Execute(c.Request.Context(), pricingusecase.CreatePriceRuleInput{
		TenantID:            tenant,
		ShipmentType:        req.ShipmentType,
		OriginOfficeID:      req.OriginOfficeID,
		DestinationOfficeID: req.DestinationOfficeID,
		Unit:                req.Unit,
		Price:               req.Price,
		Currency:            req.Currency,
		Priority:            req.Priority,
		Active:              req.Active,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": toPriceRuleResponse(*out)})
}

// Update godoc
// @Summary Actualizar regla de precios
// @Description Actualiza una regla de precios existente. Permite modificar tipo de envío, oficinas, unidad de precio, precio, moneda y prioridad. Los comodines (*) siguen siendo soportados en campos de rutas. La prioridad define el orden en búsquedas jerárquicas.
// @Tags Pricing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID de la regla" Format(uuid)
// @Param payload body PriceRuleRequest true "Solicitud de actualización con nuevos valores"
// @Success 200 {object} handler.AnyDataEnvelope "Regla de precios actualizada exitosamente"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido, payload malformado o valores inválidos"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Regla no encontrada"
// @Failure 409 {object} handler.ErrorResponse "Conflicto: nueva combinación duplicada u estado incompatible"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /pricing/rules/{id} [put]
func (h *PriceRuleHandler) Update(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	var req PriceRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	out, err := h.updateUC.Execute(c.Request.Context(), pricingusecase.UpdatePriceRuleInput{
		TenantID:            tenant,
		ID:                  id,
		ShipmentType:        req.ShipmentType,
		OriginOfficeID:      req.OriginOfficeID,
		DestinationOfficeID: req.DestinationOfficeID,
		Unit:                req.Unit,
		Price:               req.Price,
		Currency:            req.Currency,
		Priority:            req.Priority,
		Active:              req.Active,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": toPriceRuleResponse(*out)})
}

// List godoc
// @Summary Listar reglas de precios
// @Description Lista todas las reglas de precios activas del tenant actual. Incluye reglas específicas y comodines. Útil para auditoría, depuración y validación de cadenas de precios.
// @Tags Pricing
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} handler.AnyDataEnvelope "Lista de reglas de precios del tenant"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /pricing/rules [get]
func (h *PriceRuleHandler) List(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	rules, err := h.listUC.Execute(c.Request.Context(), tenant)
	if err != nil {
		_ = c.Error(err)
		return
	}

	out := make([]PriceRuleResponse, 0, len(rules))
	for _, r := range rules {
		out = append(out, toPriceRuleResponse(r))
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
}

func toPriceRuleResponse(r pricingdomain.PriceRule) PriceRuleResponse {
	return PriceRuleResponse{
		ID:                  r.ID,
		ShipmentType:        string(r.ShipmentType),
		OriginOfficeID:      r.OriginOfficeID,
		DestinationOfficeID: r.DestinationOfficeID,
		Unit:                string(r.Unit),
		Price:               r.Price,
		Currency:            r.Currency,
		Priority:            r.Priority,
		Active:              r.Active,
		CreatedAt:           r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:           r.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *PriceRuleHandler) _unused() {
	_ = time.RFC3339
}
