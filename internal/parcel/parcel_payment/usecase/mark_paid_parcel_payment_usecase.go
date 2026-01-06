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

type MarkPaidParcelPaymentUseCase struct {
	parcelRepo  coreport.ParcelReader
	paymentRepo port.ParcelPaymentRepository
	opts        coreport.TenantOptionsProvider
}

func NewMarkPaidParcelPaymentUseCase(parcelRepo coreport.ParcelReader, paymentRepo port.ParcelPaymentRepository, opts coreport.TenantOptionsProvider) *MarkPaidParcelPaymentUseCase {
	return &MarkPaidParcelPaymentUseCase{parcelRepo: parcelRepo, paymentRepo: paymentRepo, opts: opts}
}

func (u *MarkPaidParcelPaymentUseCase) Execute(ctx context.Context, tenantID string, parcelID uuid.UUID, userID *string) (*domain.ParcelPayment, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if parcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	p, err := u.parcelRepo.GetByID(ctx, tenantID, parcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": parcelID.String()}, 404)
	}

	pay, err := u.paymentRepo.GetByParcelID(ctx, tenantID, parcelID)
	if err != nil {
		return nil, err
	}
	if pay == nil {
		return nil, apperror.New("not_found", "pago no encontrado", map[string]any{"parcel_id": parcelID.String()}, 404)
	}

	if strings.TrimSpace(string(pay.Channel)) == "" {
		pay.Channel = domain.PaymentChannelCounter
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
		if o, err := u.opts.GetParcelOptions(ctx, tenantID); err == nil {
			opts = o
		} else {
			// TODO: logger
		}
	}

	isDestinationPayment := pay.PaymentType == domain.PaymentTypeFOB || pay.PaymentType == domain.PaymentTypeCollectOnDelivery

	if isDestinationPayment {
		if !opts.AllowPayInDestination {
			return nil, apperror.New("pay_in_destination_disabled", "pago en destino deshabilitado", nil, 409)
		}

		allowed := p.Status == coredomain.ParcelStatusArrivedDestination || p.Status == coredomain.ParcelStatusDelivered
		if !allowed {
			return nil, apperror.New(
				"invalid_state",
				"no se puede marcar pagado en este estado",
				map[string]any{
					"allowed": []coredomain.ParcelStatus{coredomain.ParcelStatusArrivedDestination, coredomain.ParcelStatusDelivered},
					"actual":  p.Status,
				},
				409,
			)
		}
	} else {
		allowed := p.Status == coredomain.ParcelStatusCreated || p.Status == coredomain.ParcelStatusRegistered
		if !allowed {
			return nil, apperror.New(
				"invalid_state",
				"no se puede marcar pagado en este estado",
				map[string]any{
					"allowed": []coredomain.ParcelStatus{coredomain.ParcelStatusCreated, coredomain.ParcelStatusRegistered},
					"actual":  p.Status,
				},
				409,
			)
		}
	}

	now := time.Now().UTC()
	pay.Status = domain.PaymentStatusPaid
	pay.PaidAt = &now
	pay.PaidByUserID = userID
	pay.UpdatedAt = now

	return u.paymentRepo.Upsert(ctx, tenantID, *pay)
}
