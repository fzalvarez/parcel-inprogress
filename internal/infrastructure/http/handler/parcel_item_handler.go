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
	Description string   `json:"description" binding:"required,max=200"`
	Quantity    int      `json:"quantity" binding:"required,min=1,max=9999"`
	WeightKg    float64  `json:"weight_kg" binding:"required,min=0.01,max=9999"`
	LengthCm    *float64 `json:"length_cm" binding:"omitempty,min=0.01,max=9999"`
	WidthCm     *float64 `json:"width_cm" binding:"omitempty,min=0.01,max=9999"`
	HeightCm    *float64 `json:"height_cm" binding:"omitempty,min=0.01,max=9999"`
	UnitPrice   float64  `json:"unit_price" binding:"omitempty,min=0,max=999999"`
	ContentType *string  `json:"content_type" binding:"omitempty,max=100"`
	Notes       *string  `json:"notes" binding:"omitempty,max=300"`
}

type ParcelItemResponse struct {
	ID               string   `json:"id"`
	ParcelID         string   `json:"parcel_id"`
	Description      string   `json:"description"`
	Quantity         int      `json:"quantity"`
	WeightKg         float64  `json:"weight_kg"`
	LengthCm         *float64 `json:"length_cm,omitempty"`
	WidthCm          *float64 `json:"width_cm,omitempty"`
	HeightCm         *float64 `json:"height_cm,omitempty"`
	VolumetricWeight *float64 `json:"volumetric_weight,omitempty"`
	BillableWeight   float64  `json:"billable_weight"`
	UnitPrice        float64  `json:"unit_price"`
	ContentType      *string  `json:"content_type,omitempty"`
	Notes            *string  `json:"notes,omitempty"`
	CreatedAt        string   `json:"created_at"`
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
// @Summary Agregar artículo (item) al envío
// @Description Agrega un bulto/artículo al envío con cálculo automático de peso facturable y precio. Soporta dimensiones opcionales (largo, ancho, alto) para cálculo de peso volumétrico. El peso facturable se calcula como máximo entre peso real y volumétrico (si aplica según configuración del tenant). El precio unitario se busca mediante reglas de precios jerárquicas.
// @Tags ParcelItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Param payload body CreateParcelItemRequest true "Datos del item (descripción y peso requeridos, dimensiones opcionales para volumétrico)"
// @Success 201 {object} handler.AnyDataEnvelope "Item creado exitosamente con peso y precio calculados"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido, payload malformado o valores inválidos"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío no encontrado"
// @Failure 409 {object} handler.ErrorResponse "Conflicto: regla de precios no encontrada (sugerencia: use comodines *) o estado del envío no permite agregar items"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
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

	item, err := h.addUC.Execute(c.Request.Context(), itemusecase.AddParcelItemInput{
		TenantID:    tenant,
		UserID:      strings.TrimSpace(anyToString(userID)),
		UserName:    strings.TrimSpace(anyToString(userName)),
		ParcelID:    parcelID,
		Description: req.Description,
		Quantity:    req.Quantity,
		WeightKg:    req.WeightKg,
		LengthCm:    req.LengthCm,
		WidthCm:     req.WidthCm,
		HeightCm:    req.HeightCm,
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
			ID:               item.ID,
			ParcelID:         item.ParcelID,
			Description:      item.Description,
			Quantity:         item.Quantity,
			WeightKg:         item.WeightKg,
			LengthCm:         item.LengthCm,
			WidthCm:          item.WidthCm,
			HeightCm:         item.HeightCm,
			VolumetricWeight: item.VolumetricWeight,
			BillableWeight:   item.BillableWeight,
			UnitPrice:        item.UnitPrice,
			ContentType:      item.ContentType,
			Notes:            item.Notes,
			CreatedAt:        item.CreatedAt.UTC().Format(time.RFC3339),
		},
	})
}

// List godoc
// @Summary Listar artículos del envío
// @Description Lista todos los artículos (items/bultos) agregados a un envío. Devuelve detalles de cada item incluyendo dimensiones, pesos calculados (real, volumétrico, facturable) y precio unitario.
// @Tags ParcelItems
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope "Lista de items del envío"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío no encontrado"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
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
			ID:               it.ID,
			ParcelID:         it.ParcelID,
			Description:      it.Description,
			Quantity:         it.Quantity,
			WeightKg:         it.WeightKg,
			LengthCm:         it.LengthCm,
			WidthCm:          it.WidthCm,
			HeightCm:         it.HeightCm,
			VolumetricWeight: it.VolumetricWeight,
			BillableWeight:   it.BillableWeight,
			UnitPrice:        it.UnitPrice,
			ContentType:      it.ContentType,
			Notes:            it.Notes,
			CreatedAt:        it.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
}

// Delete godoc
// @Summary Eliminar artículo del envío
// @Description Elimina un artículo específico agregado a un envío. Solo permitido en ciertos estados del envío (antes de registro o bajo condiciones especiales).
// @Tags ParcelItems
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Param item_id path string true "UUID del item a eliminar" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope "Item eliminado exitosamente"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id o item_id inválidos"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío o item no encontrado"
// @Failure 409 {object} handler.ErrorResponse "Conflicto: no se puede eliminar item en este estado del envío"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
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
