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
	PaymentType  string  `json:"payment_type" binding:"required,oneof=CASH FOB CARD TRANSFER EWALLET FREE COLLECT_ON_DELIVERY"`
	Currency     *string `json:"currency" binding:"omitempty,oneof=PEN USD"`
	Amount       float64 `json:"amount" binding:"required"`
	Notes        *string `json:"notes" binding:"omitempty,max=200"`
	Channel      *string `json:"channel" binding:"omitempty,oneof=COUNTER WEB"`
	OfficeID     *string `json:"office_id" binding:"omitempty,uuid"`
	CashboxID    *string `json:"cashbox_id" binding:"omitempty,max=50"`
	SellerUserID *string `json:"seller_user_id" binding:"omitempty,uuid"`
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
	Channel      string  `json:"channel"`
	OfficeID     *string `json:"office_id,omitempty"`
	CashboxID    *string `json:"cashbox_id,omitempty"`
	SellerUserID *string `json:"seller_user_id,omitempty"`
}

type ParcelPaymentHandler struct {
	upsertUC   *paymentusecase.UpsertParcelPaymentUseCase
	getUC      *paymentusecase.GetParcelPaymentUseCase
	markPaidUC *paymentusecase.MarkPaidParcelPaymentUseCase
}

func NewParcelPaymentHandler(upsertUC *paymentusecase.UpsertParcelPaymentUseCase, getUC *paymentusecase.GetParcelPaymentUseCase, markPaidUC *paymentusecase.MarkPaidParcelPaymentUseCase) *ParcelPaymentHandler {
	return &ParcelPaymentHandler{upsertUC: upsertUC, getUC: getUC, markPaidUC: markPaidUC}
}

// Upsert godoc
// @Summary Crear o actualizar información de pago
// @Description Crea o actualiza la información de pago del envío. Incluye tipo de pago, monto, moneda, canal, oficina y datos de caja. Soporta múltiples formas de pago (CASH, FOB, CARD, TRANSFER, EWALLET, FREE, COLLECT_ON_DELIVERY). El estado inicial es PENDING.
// @Tags ParcelPayments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Param payload body UpsertParcelPaymentRequest true "Solicitud con datos de pago (monto requerido, type de pago, moneda por defecto PEN)"
// @Success 200 {object} handler.AnyDataEnvelope "Pago creado o actualizado exitosamente"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido, payload malformado o valores inválidos"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío no encontrado"
// @Failure 409 {object} handler.ErrorResponse "Conflicto: estado incompatible o envío no permite esta operación"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /parcels/{id}/payment [put]
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

	channel := "COUNTER"
	if req.Channel != nil && strings.TrimSpace(*req.Channel) != "" {
		channel = strings.TrimSpace(*req.Channel)
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	pay, err := h.upsertUC.Execute(c.Request.Context(), paymentusecase.UpsertParcelPaymentInput{
		TenantID:     tenant,
		ParcelID:     parcelID,
		PaymentType:  paymentdomain.PaymentType(strings.TrimSpace(req.PaymentType)),
		Currency:     paymentdomain.Currency(currency),
		Amount:       req.Amount,
		Notes:        req.Notes,
		Channel:      paymentdomain.PaymentChannel(channel),
		OfficeID:     req.OfficeID,
		CashboxID:    req.CashboxID,
		SellerUserID: req.SellerUserID,
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
			Channel:      string(pay.Channel),
			OfficeID:     pay.OfficeID,
			CashboxID:    pay.CashboxID,
			SellerUserID: pay.SellerUserID,
		},
	})
}

// Get godoc
// @Summary Obtener información de pago del envío
// @Description Devuelve los detalles completos del pago registrado para un envío. Incluye tipo de pago, monto, estado (PENDING/PAID), moneda, canal, oficina y datos de caja.
// @Tags ParcelPayments
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope "Información de pago obtenida"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío o pago no encontrado"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /parcels/{id}/payment [get]
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
			Channel:      string(pay.Channel),
			OfficeID:     pay.OfficeID,
			CashboxID:    pay.CashboxID,
			SellerUserID: pay.SellerUserID,
		},
	})
}

// MarkPaid godoc
// @Summary Marcar pago como realizado
// @Description Transiciona el pago a estado PAID. Intégrase con el servicio de caja (CASHBOX) para registrar la transacción y confirmar la recaudación. Captura el user_id del operador que marca el pago.
// @Tags ParcelPayments
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope "Pago marcado como realizado exitosamente"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío o pago no encontrado"
// @Failure 409 {object} handler.ErrorResponse "Conflicto: pago ya realizado o estado no permitido"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor o fallo de integración con caja"
// @Router /parcels/{id}/payment/mark-paid [post]
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
			Channel:      string(pay.Channel),
			OfficeID:     pay.OfficeID,
			CashboxID:    pay.CashboxID,
			SellerUserID: pay.SellerUserID,
		},
	})
}
