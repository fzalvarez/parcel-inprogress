package port

import (
	"context"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
)

type ParcelReader interface {
	ListByFilters(ctx context.Context, tenantID string, f coreport.ListParcelFilters) ([]domain.Parcel, error)
}
