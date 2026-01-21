package postgres

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	trackingdomain "ms-parcel-core/internal/parcel/parcel_tracking/domain"
)

// DBTrackingEvent representa el modelo de base de datos para TrackingEvent
type DBTrackingEvent struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ParcelID   string    `gorm:"type:varchar(100);not null;index"`
	EventType  string    `gorm:"type:varchar(50);not null"`
	OccurredAt time.Time `gorm:"not null;index"`
	UserID     string    `gorm:"type:varchar(100);not null"`
	UserName   string    `gorm:"type:varchar(255)"`
	Metadata   *string   `gorm:"type:jsonb"`
}

func (DBTrackingEvent) TableName() string {
	return "tracking_events"
}

// ToDomain convierte DBTrackingEvent a trackingdomain.TrackingEvent
func (db *DBTrackingEvent) ToDomain() trackingdomain.TrackingEvent {
	var metadata map[string]any
	if db.Metadata != nil && *db.Metadata != "" {
		_ = json.Unmarshal([]byte(*db.Metadata), &metadata)
	}

	id, _ := uuid.Parse(db.ID.String())
	return trackingdomain.TrackingEvent{
		ID:         id,
		ParcelID:   db.ParcelID,
		EventType:  db.EventType,
		OccurredAt: db.OccurredAt,
		UserID:     db.UserID,
		UserName:   db.UserName,
		Metadata:   metadata,
	}
}

// FromDomain convierte trackingdomain.TrackingEvent a DBTrackingEvent
func (db *DBTrackingEvent) FromDomain(evt trackingdomain.TrackingEvent) error {
	var metadataJSON *string
	if evt.Metadata != nil {
		data, err := json.Marshal(evt.Metadata)
		if err == nil {
			str := string(data)
			metadataJSON = &str
		}
	}

	*db = DBTrackingEvent{
		ID:         evt.ID,
		ParcelID:   evt.ParcelID,
		EventType:  evt.EventType,
		OccurredAt: evt.OccurredAt,
		UserID:     evt.UserID,
		UserName:   evt.UserName,
		Metadata:   metadataJSON,
	}
	return nil
}

// BeforeCreate hook de GORM
func (db *DBTrackingEvent) BeforeCreate(tx *gorm.DB) error {
	if db.ID == uuid.Nil {
		db.ID = uuid.New()
	}
	return nil
}
