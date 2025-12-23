package domain

import (
	"time"

	"github.com/google/uuid"
)

type TrackingEvent struct {
	ID         uuid.UUID
	ParcelID   string
	EventType  string
	OccurredAt time.Time
	UserID     string
	UserName   string
	Metadata   map[string]any
}
