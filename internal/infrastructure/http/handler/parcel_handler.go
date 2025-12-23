package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_core/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type CreateParcelRequest struct {
	ShipmentType        string  `json:"shipment_type" binding:"required,oneof=BUS CARGUERO"`
	OriginOfficeID      string  `json:"origin_office_id" binding:"required,uuid"`
	DestinationOfficeID string  `json:"destination_office_id" binding:"required,uuid"`
	SenderPersonID      string  `json:"sender_person_id" binding:"required,uuid"`
	RecipientPersonID   string  `json:"recipient_person_id" binding:"required,uuid"`
	Notes               *string `json:"notes" binding:"omitempty,max=500"`
	PackageKey          string  `json:"package_key" binding:"required,min=4,max=200"`
	PackageKeyConfirm   string  `json:"package_key_confirm" binding:"required,min=4,max=200"`
}

type CreateParcelResponse struct {
	ID                  string  `json:"id"`
	Status              string  `json:"status"`
	ShipmentType        string  `json:"shipment_type"`
	OriginOfficeID      string  `json:"origin_office_id"`
	DestinationOfficeID string  `json:"destination_office_id"`
	SenderPersonID      string  `json:"sender_person_id"`
	RecipientPersonID   string  `json:"recipient_person_id"`
	Notes               *string `json:"notes"`
	CreatedAt           string  `json:"created_at"`
}

type ParcelHandler struct {
	createUC *usecase.CreateParcelUseCase
}

func NewParcelHandler(createUC *usecase.CreateParcelUseCase) *ParcelHandler {
	return &ParcelHandler{createUC: createUC}
}

func (h *ParcelHandler) Create(c *gin.Context) {
	var req CreateParcelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	in := usecase.CreateParcelInput{
		TenantID:            strings.TrimSpace(anyToString(tenantID)),
		UserID:              strings.TrimSpace(anyToString(userID)),
		UserName:            strings.TrimSpace(anyToString(userName)),
		ShipmentType:        domain.ShipmentType(req.ShipmentType),
		OriginOfficeID:      req.OriginOfficeID,
		DestinationOfficeID: req.DestinationOfficeID,
		SenderPersonID:      req.SenderPersonID,
		RecipientPersonID:   req.RecipientPersonID,
		Notes:               req.Notes,
		PackageKey:          req.PackageKey,
		PackageKeyConfirm:   req.PackageKeyConfirm,
	}

	id, err := h.createUC.Execute(c.Request.Context(), in)
	if err != nil {
		_ = c.Error(err)
		return
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)

	// Header temporal de debug (remover luego de estabilizar)
	c.Header("X-Debug-Parcel", "created")

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": CreateParcelResponse{
			ID:                  id.String(),
			Status:              string(domain.ParcelStatusCreated),
			ShipmentType:        req.ShipmentType,
			OriginOfficeID:      req.OriginOfficeID,
			DestinationOfficeID: req.DestinationOfficeID,
			SenderPersonID:      req.SenderPersonID,
			RecipientPersonID:   req.RecipientPersonID,
			Notes:               req.Notes,
			CreatedAt:           createdAt,
		},
	})
}

func anyToString(v any) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return ""
}
