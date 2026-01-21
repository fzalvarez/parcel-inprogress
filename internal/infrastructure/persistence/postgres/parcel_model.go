package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
)

// DBParcel representa el modelo de base de datos para Parcel
type DBParcel struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TenantID             string    `gorm:"type:varchar(100);not null;index:idx_parcel_tenant"`
	OriginOfficeID       string    `gorm:"type:varchar(100);not null;index"`
	DestinationOfficeID  string    `gorm:"type:varchar(100);not null;index"`
	SenderPersonID       string    `gorm:"type:varchar(100);not null"`
	RecipientPersonID    string    `gorm:"type:varchar(100);not null"`
	ShipmentType         string    `gorm:"type:varchar(50);not null"`
	Notes                *string   `gorm:"type:text"`
	PackageKeyHashSHA256 string    `gorm:"type:varchar(255)"`
	Status               string    `gorm:"type:varchar(50);not null;index"`
	TrackingCode         string    `gorm:"type:varchar(50);uniqueIndex"`
	CreatedByUserID      string    `gorm:"type:varchar(100);not null"`
	CreatedByUserName    string    `gorm:"type:varchar(255)"`
	CreatedAt            time.Time `gorm:"not null"`
	RegisteredAt         *time.Time

	BoardedVehicleID   *string `gorm:"type:varchar(100)"`
	BoardedTripID      *string `gorm:"type:varchar(100)"`
	BoardedDepartureAt *time.Time
	BoardedAt          *time.Time
	BoardedByUserID    *string `gorm:"type:varchar(100)"`
	DeliveredAt        *time.Time
	DeliveredByUserID  *string `gorm:"type:varchar(100)"`
	ArrivedAt          *time.Time
	ArrivedByUserID    *string `gorm:"type:varchar(100)"`
	DepartedAt         *time.Time
	DepartedByUserID   *string `gorm:"type:varchar(100)"`
}

func (DBParcel) TableName() string {
	return "parcels"
}

// ToDomain convierte DBParcel a domain.Parcel
func (db *DBParcel) ToDomain() domain.Parcel {
	return domain.Parcel{
		ID:                   db.ID.String(),
		TenantID:             db.TenantID,
		OriginOfficeID:       db.OriginOfficeID,
		DestinationOfficeID:  db.DestinationOfficeID,
		SenderPersonID:       db.SenderPersonID,
		RecipientPersonID:    db.RecipientPersonID,
		ShipmentType:         domain.ShipmentType(db.ShipmentType),
		Notes:                db.Notes,
		PackageKeyHashSHA256: db.PackageKeyHashSHA256,
		Status:               domain.ParcelStatus(db.Status),
		TrackingCode:         db.TrackingCode,
		CreatedByUserID:      db.CreatedByUserID,
		CreatedByUserName:    db.CreatedByUserName,
		CreatedAt:            db.CreatedAt,
		RegisteredAt:         db.RegisteredAt,
		BoardedVehicleID:     db.BoardedVehicleID,
		BoardedTripID:        db.BoardedTripID,
		BoardedDepartureAt:   db.BoardedDepartureAt,
		BoardedAt:            db.BoardedAt,
		BoardedByUserID:      db.BoardedByUserID,
		DeliveredAt:          db.DeliveredAt,
		DeliveredByUserID:    db.DeliveredByUserID,
		ArrivedAt:            db.ArrivedAt,
		ArrivedByUserID:      db.ArrivedByUserID,
		DepartedAt:           db.DepartedAt,
		DepartedByUserID:     db.DepartedByUserID,
	}
}

// FromDomain convierte domain.Parcel a DBParcel
func (db *DBParcel) FromDomain(p domain.Parcel) error {
	id, err := uuid.Parse(p.ID)
	if err != nil && p.ID != "" {
		return err
	}
	if p.ID == "" {
		id = uuid.New()
	}

	*db = DBParcel{
		ID:                   id,
		TenantID:             p.TenantID,
		OriginOfficeID:       p.OriginOfficeID,
		DestinationOfficeID:  p.DestinationOfficeID,
		SenderPersonID:       p.SenderPersonID,
		RecipientPersonID:    p.RecipientPersonID,
		ShipmentType:         string(p.ShipmentType),
		Notes:                p.Notes,
		PackageKeyHashSHA256: p.PackageKeyHashSHA256,
		Status:               string(p.Status),
		TrackingCode:         p.TrackingCode,
		CreatedByUserID:      p.CreatedByUserID,
		CreatedByUserName:    p.CreatedByUserName,
		CreatedAt:            p.CreatedAt,
		RegisteredAt:         p.RegisteredAt,
		BoardedVehicleID:     p.BoardedVehicleID,
		BoardedTripID:        p.BoardedTripID,
		BoardedDepartureAt:   p.BoardedDepartureAt,
		BoardedAt:            p.BoardedAt,
		BoardedByUserID:      p.BoardedByUserID,
		DeliveredAt:          p.DeliveredAt,
		DeliveredByUserID:    p.DeliveredByUserID,
		ArrivedAt:            p.ArrivedAt,
		ArrivedByUserID:      p.ArrivedByUserID,
		DepartedAt:           p.DepartedAt,
		DepartedByUserID:     p.DepartedByUserID,
	}
	return nil
}

// BeforeCreate hook de GORM
func (db *DBParcel) BeforeCreate(tx *gorm.DB) error {
	if db.ID == uuid.Nil {
		db.ID = uuid.New()
	}
	return nil
}
