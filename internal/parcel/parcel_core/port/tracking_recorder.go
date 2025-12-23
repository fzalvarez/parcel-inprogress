package port

import (
	"context"
	"time"
)

type TrackingEventDTO struct {
	ParcelID   string
	EventType  string
	OccurredAt time.Time
	UserID     string
	UserName   string
	Metadata   map[string]any
}

const (
	EventTypeParcelCreated            = "PARCEL_CREATED"
	EventTypeParcelRegistered         = "PARCEL_REGISTERED"
	EventTypeParcelBoarded            = "PARCEL_BOARDED"
	EventTypeParcelInTransit          = "PARCEL_IN_TRANSIT"
	EventTypeParcelArrivedDestination = "PARCEL_ARRIVED_DESTINATION"
	EventTypeParcelDelivered          = "PARCEL_DELIVERED"
)

type TrackingRecorder interface {
	RecordEvent(ctx context.Context, tenantID string, ev TrackingEventDTO) error
}
