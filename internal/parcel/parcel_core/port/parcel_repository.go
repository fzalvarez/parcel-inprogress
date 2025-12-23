package port

import (
	"context"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
)

type ParcelRepository interface {
	Create(ctx context.Context, p domain.Parcel) (uuid.UUID, error)
}
