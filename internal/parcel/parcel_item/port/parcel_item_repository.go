package port

import (
	"context"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_item/domain"
)

type ParcelItemRepository interface {
	Add(ctx context.Context, tenantID string, item domain.ParcelItem) (uuid.UUID, error)
	ListByParcelID(ctx context.Context, tenantID string, parcelID uuid.UUID) ([]domain.ParcelItem, error)
	Delete(ctx context.Context, tenantID string, parcelID uuid.UUID, itemID uuid.UUID) error
}
