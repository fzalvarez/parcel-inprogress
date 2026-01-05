package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	coredomain "ms-parcel-core/internal/parcel/parcel_core/domain"
	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_payment/domain"
	"ms-parcel-core/internal/parcel/parcel_payment/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type UpsertParcelPaymentInput struct {
	TenantID    string
	ParcelID    uuid.UUID
	PaymentType domain.PaymentType
	Currency    domain.Currency
	Amount      float64
	Notes       *string
}

type UpsertParcelPaymentUseCase struct {
	parcelRepo  coreport.ParcelReader
	paymentRepo port.ParcelPaymentRepository
	opts        coreport.TenantOptionsProvider
}

func NewUpsertParcelPaymentUseCase(parcelRepo coreport.ParcelReader, paymentRepo port.ParcelPaymentRepository, opts coreport.TenantOptionsProvider) *UpsertParcelPaymentUseCase {
	return &UpsertParcelPaymentUseCase{parcelRepo: parcelRepo, paymentRepo: paymentRepo, opts: opts}
}

func (u *UpsertParcelPaymentUseCase) Execute(ctx context.Context, in UpsertParcelPaymentInput) (*domain.ParcelPayment, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}
	if in.Amount < 0 {
		return nil, apperror.NewBadRequest("validation_error", "amount inválido", map[string]any{"field": "amount"})
	}
	if in.PaymentType != domain.PaymentTypeFree && in.Amount <= 0 {
		return nil, apperror.NewBadRequest("validation_error", "amount debe ser > 0", map[string]any{"field": "amount"})
	}

	p, err := u.parcelRepo.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	allowed := p.Status == coredomain.ParcelStatusCreated || p.Status == coredomain.ParcelStatusRegistered
	if !allowed {
		return nil, apperror.New("invalid_state", "no se puede modificar el pago en este estado", map[string]any{"allowed": []coredomain.ParcelStatus{coredomain.ParcelStatusCreated, coredomain.ParcelStatusRegistered}, "actual": p.Status}, 409)
	}

	defaults := coreport.ParcelOptions{
		RequirePackageKey:       true,
		UsePriceTable:           true,
		AllowManualPrice:        false,
		AllowOverridePriceTable: true,
		AllowPayInDestination:   false,
	}
	opts := defaults
	if u.opts != nil {
		if o, err := u.opts.GetParcelOptions(ctx, in.TenantID); err == nil {
			opts = o
		} else {
			// TODO: logger
		}
	}

	if !opts.AllowPayInDestination {
		if in.PaymentType == domain.PaymentTypeFOB || in.PaymentType == domain.PaymentTypeCollectOnDelivery {
			return nil, apperror.New("pay_in_destination_disabled", "pago en destino deshabilitado", nil, 409)
		}
	}

	// TODO: no calcular monto desde items aún.

	now := time.Now().UTC()
	existing, err := u.paymentRepo.GetByParcelID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}

	pay := domain.ParcelPayment{
		ID:          uuid.NewString(),
		TenantID:    in.TenantID,
		ParcelID:    in.ParcelID.String(),
		PaymentType: in.PaymentType,
		Currency:    in.Currency,
		Amount:      in.Amount,
		Notes:       in.Notes,
		Status:      domain.PaymentStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if existing != nil {
		pay.ID = existing.ID
		pay.CreatedAt = existing.CreatedAt
	}

	return u.paymentRepo.Upsert(ctx, in.TenantID, pay)
}
