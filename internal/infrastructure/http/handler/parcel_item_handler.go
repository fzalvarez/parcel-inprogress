package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	itemusecase "ms-parcel-core/internal/parcel/parcel_item/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type CreateParcelItemRequest struct {
	Description string  `json:"description" binding:"required,max=200"`
	Quantity    int     `json:"quantity" binding:"required,min=1,max=9999"`
	WeightKg    float64 `json:"weight_kg" binding:"required,min=0.01,max=9999"`
	UnitPrice   float64 `json:"unit_price" binding:"required,min=0,max=999999"`
	ContentType *string `json:"content_type" binding:"omitempty,max=100"`
	Notes       *string `json:"notes" binding:"omitempty,max=300"`
}

type ParcelItemResponse struct {
	ID          string  `json:"id"`
	ParcelID    string  `json:"parcel_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	WeightKg    float64 `json:"weight_kg"`
	UnitPrice   float64 `json:"unit_price"`
	ContentType *string `json:"content_type,omitempty"`
	Notes       *string `json:"notes,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type ParcelItemHandler struct {
	addUC    *itemusecase.AddParcelItemUseCase
	listUC   *itemusecase.ListParcelItemsUseCase
	deleteUC *itemusecase.DeleteParcelItemUseCase
}

func NewParcelItemHandler(addUC *itemusecase.AddParcelItemUseCase, listUC *itemusecase.ListParcelItemsUseCase, deleteUC *itemusecase.DeleteParcelItemUseCase) *ParcelItemHandler {
	return &ParcelItemHandler{addUC: addUC, listUC: listUC, deleteUC: deleteUC}
}

// Add godoc
// @Summary Agregar item
// @Description Agrega un item al envío
// @Tags ParcelItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Param payload body AddParcelItemRequest true "Add item"
// @Success 200 {object} handler.AnyDataEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/items [post]
func (h *ParcelItemHandler) Add(c *gin.Context) {
	parcelIDStr := strings.TrimSpace(c.Param("id"))
	parcelID, err := uuid.Parse(parcelIDStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	var req CreateParcelItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	id, err := h.addUC.Execute(c.Request.Context(), itemusecase.AddParcelItemInput{
		TenantID:    tenant,
		UserID:      strings.TrimSpace(anyToString(userID)),
		UserName:    strings.TrimSpace(anyToString(userName)),
		ParcelID:    parcelID,
		Description: req.Description,
		Quantity:    req.Quantity,
		WeightKg:    req.WeightKg,
		UnitPrice:   req.UnitPrice,
		ContentType: req.ContentType,
		Notes:       req.Notes,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": ParcelItemResponse{
			ID:          id.String(),
			ParcelID:    parcelID.String(),
			Description: req.Description,
			Quantity:    req.Quantity,
			WeightKg:    req.WeightKg,
			UnitPrice:   req.UnitPrice,
			ContentType: req.ContentType,
			Notes:       req.Notes,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// List godoc
// @Summary Listar items
// @Description Lista items del envío
// @Tags ParcelItems
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/items [get]
func (h *ParcelItemHandler) List(c *gin.Context) {
	parcelIDStr := strings.TrimSpace(c.Param("id"))
	parcelID, err := uuid.Parse(parcelIDStr)
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

	items, err := h.listUC.Execute(c.Request.Context(), tenant, parcelID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	out := make([]ParcelItemResponse, 0, len(items))
	for _, it := range items {
		out = append(out, ParcelItemResponse{
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

	c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
}

// Delete godoc
// @Summary Eliminar item
// @Description Elimina un item del envío
// @Tags ParcelItems
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Param item_id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/items/{item_id} [delete]
func (h *ParcelItemHandler) Delete(c *gin.Context) {
	parcelIDStr := strings.TrimSpace(c.Param("id"))
	parcelID, err := uuid.Parse(parcelIDStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	itemIDStr := strings.TrimSpace(c.Param("item_id"))
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "item_id inválido", map[string]any{"field": "item_id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	if err := h.deleteUC.Execute(c.Request.Context(), itemusecase.DeleteParcelItemInput{
		TenantID: tenant,
		UserID:   strings.TrimSpace(anyToString(userID)),
		UserName: strings.TrimSpace(anyToString(userName)),
		ParcelID: parcelID,
		ItemID:   itemID,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"deleted": true}})
}
