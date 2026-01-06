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
// @Description Crea una regla de precios por tenant
// @Tags Pricing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param payload body PriceRuleRequest true "Price rule"
// @Success 200 {object} handler.AnyDataEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
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
// @Description Actualiza una regla de precios
// @Tags Pricing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Param payload body PriceRuleRequest true "Price rule"
// @Success 200 {object} handler.AnyDataEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
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
// @Description Lista reglas de precios por tenant
// @Tags Pricing
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} handler.AnyDataEnvelope
// @Failure 401 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
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
		Active:              r.Active,
		CreatedAt:           r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:           r.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *PriceRuleHandler) _unused() {
	_ = time.RFC3339
}
