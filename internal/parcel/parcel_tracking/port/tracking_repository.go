package port

import (
	"context"

	"ms-parcel-core/internal/parcel/parcel_tracking/domain"
)

type TrackingRepository interface {
	Append(ctx context.Context, tenantID string, ev domain.TrackingEvent) error
	ListByParcelID(ctx context.Context, tenantID string, parcelID string) ([]domain.TrackingEvent, error)
}
