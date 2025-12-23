package port

import (
	"context"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
)

type ParcelReader interface {
	GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*domain.Parcel, error)
}
