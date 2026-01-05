package port

import (
	"context"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_payment/domain"
)

type ParcelPaymentRepository interface {
	Upsert(ctx context.Context, tenantID string, p domain.ParcelPayment) (*domain.ParcelPayment, error)
	GetByParcelID(ctx context.Context, tenantID string, parcelID uuid.UUID) (*domain.ParcelPayment, error)
}
