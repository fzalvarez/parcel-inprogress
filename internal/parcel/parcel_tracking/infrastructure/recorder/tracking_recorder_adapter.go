package recorder

import (
	"context"

	"github.com/google/uuid"

	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_tracking/domain"
	trackingport "ms-parcel-core/internal/parcel/parcel_tracking/port"
)

type TrackingRecorderAdapter struct {
	repo trackingport.TrackingRepository
}

var _ coreport.TrackingRecorder = (*TrackingRecorderAdapter)(nil)

func NewTrackingRecorderAdapter(repo trackingport.TrackingRepository) *TrackingRecorderAdapter {
	return &TrackingRecorderAdapter{repo: repo}
}

func (a *TrackingRecorderAdapter) RecordEvent(ctx context.Context, tenantID string, ev coreport.TrackingEventDTO) error {
	te := domain.TrackingEvent{
		ID:         uuid.New(),
		ParcelID:   ev.ParcelID,
		EventType:  ev.EventType,
		OccurredAt: ev.OccurredAt,
		UserID:     ev.UserID,
		UserName:   ev.UserName,
		Metadata:   ev.Metadata,
	}
	return a.repo.Append(ctx, tenantID, te)
}
