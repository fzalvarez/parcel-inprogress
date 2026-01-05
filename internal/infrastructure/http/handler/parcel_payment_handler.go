package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	paymentdomain "ms-parcel-core/internal/parcel/parcel_payment/domain"
	paymentusecase "ms-parcel-core/internal/parcel/parcel_payment/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type UpsertParcelPaymentRequest struct {
	PaymentType string  `json:"payment_type" binding:"required,oneof=CASH FOB CARD TRANSFER EWALLET FREE COLLECT_ON_DELIVERY"`
	Currency    *string `json:"currency" binding:"omitempty,oneof=PEN USD"`
	Amount      float64 `json:"amount" binding:"required"`
	Notes       *string `json:"notes" binding:"omitempty,max=200"`
}

type ParcelPaymentResponse struct {
	ID           string  `json:"id"`
	ParcelID     string  `json:"parcel_id"`
	PaymentType  string  `json:"payment_type"`
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount"`
	Notes        *string `json:"notes,omitempty"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	PaidAt       *string `json:"paid_at,omitempty"`
	PaidByUserID *string `json:"paid_by_user_id,omitempty"`
}

type ParcelPaymentHandler struct {
	upsertUC   *paymentusecase.UpsertParcelPaymentUseCase
	getUC      *paymentusecase.GetParcelPaymentUseCase
	markPaidUC *paymentusecase.MarkPaidParcelPaymentUseCase
}

func NewParcelPaymentHandler(upsertUC *paymentusecase.UpsertParcelPaymentUseCase, getUC *paymentusecase.GetParcelPaymentUseCase, markPaidUC *paymentusecase.MarkPaidParcelPaymentUseCase) *ParcelPaymentHandler {
	return &ParcelPaymentHandler{upsertUC: upsertUC, getUC: getUC, markPaidUC: markPaidUC}
}

func (h *ParcelPaymentHandler) Upsert(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	parcelID, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	var req UpsertParcelPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	currency := "PEN"
	if req.Currency != nil && strings.TrimSpace(*req.Currency) != "" {
		currency = strings.TrimSpace(*req.Currency)
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	pay, err := h.upsertUC.Execute(c.Request.Context(), paymentusecase.UpsertParcelPaymentInput{
		TenantID:    tenant,
		ParcelID:    parcelID,
		PaymentType: paymentdomain.PaymentType(strings.TrimSpace(req.PaymentType)),
		Currency:    paymentdomain.Currency(currency),
		Amount:      req.Amount,
		Notes:       req.Notes,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var paidAtStr *string
	if pay.PaidAt != nil {
		s := pay.PaidAt.UTC().Format(time.RFC3339)
		paidAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": ParcelPaymentResponse{
			ID:           pay.ID,
			ParcelID:     pay.ParcelID,
			PaymentType:  string(pay.PaymentType),
			Currency:     string(pay.Currency),
			Amount:       pay.Amount,
			Notes:        pay.Notes,
			Status:       string(pay.Status),
			CreatedAt:    pay.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:    pay.UpdatedAt.UTC().Format(time.RFC3339),
			PaidAt:       paidAtStr,
			PaidByUserID: pay.PaidByUserID,
		},
	})
}

func (h *ParcelPaymentHandler) Get(c *gin.Context) {
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

	pay, err := h.getUC.Execute(c.Request.Context(), tenant, parcelID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if pay == nil {
		_ = c.Error(apperror.New("not_found", "pago no encontrado", map[string]any{"parcel_id": parcelID.String()}, 404))
		return
	}

	var paidAtStr *string
	if pay.PaidAt != nil {
		s := pay.PaidAt.UTC().Format(time.RFC3339)
		paidAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": ParcelPaymentResponse{
			ID:           pay.ID,
			ParcelID:     pay.ParcelID,
			PaymentType:  string(pay.PaymentType),
			Currency:     string(pay.Currency),
			Amount:       pay.Amount,
			Notes:        pay.Notes,
			Status:       string(pay.Status),
			CreatedAt:    pay.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:    pay.UpdatedAt.UTC().Format(time.RFC3339),
			PaidAt:       paidAtStr,
			PaidByUserID: pay.PaidByUserID,
		},
	})
}

func (h *ParcelPaymentHandler) MarkPaid(c *gin.Context) {
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

	userIDVal, _ := c.Get("user_id")
	uid := strings.TrimSpace(anyToString(userIDVal))
	var uidPtr *string
	if uid != "" {
		uidPtr = &uid
	}

	pay, err := h.markPaidUC.Execute(c.Request.Context(), tenant, parcelID, uidPtr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var paidAtStr *string
	if pay.PaidAt != nil {
		s := pay.PaidAt.UTC().Format(time.RFC3339)
		paidAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": ParcelPaymentResponse{
			ID:           pay.ID,
			ParcelID:     pay.ParcelID,
			PaymentType:  string(pay.PaymentType),
			Currency:     string(pay.Currency),
			Amount:       pay.Amount,
			Notes:        pay.Notes,
			Status:       string(pay.Status),
			CreatedAt:    pay.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:    pay.UpdatedAt.UTC().Format(time.RFC3339),
			PaidAt:       paidAtStr,
			PaidByUserID: pay.PaidByUserID,
		},
	})
}
