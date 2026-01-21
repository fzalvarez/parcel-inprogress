package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ms-parcel-core/internal/parcel/parcel_item/domain"
)

// DBParcelItem representa el modelo de base de datos para ParcelItem
type DBParcelItem struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ParcelID         uuid.UUID `gorm:"type:uuid;not null;index"`
	Description      string    `gorm:"type:varchar(500);not null"`
	Quantity         int       `gorm:"not null"`
	WeightKg         float64   `gorm:"type:decimal(10,2);not null"`
	LengthCm         *float64  `gorm:"type:decimal(10,2)"`
	WidthCm          *float64  `gorm:"type:decimal(10,2)"`
	HeightCm         *float64  `gorm:"type:decimal(10,2)"`
	VolumetricWeight *float64  `gorm:"type:decimal(10,2)"`
	BillableWeight   float64   `gorm:"type:decimal(10,2);not null"`
	UnitPrice        float64   `gorm:"type:decimal(10,2);not null"`
	ContentType      *string   `gorm:"type:varchar(100)"`
	Notes            *string   `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"not null"`
}

func (DBParcelItem) TableName() string {
	return "parcel_items"
}

// ToDomain convierte DBParcelItem a domain.ParcelItem
func (db *DBParcelItem) ToDomain() domain.ParcelItem {
	return domain.ParcelItem{
		ID:               db.ID.String(),
		ParcelID:         db.ParcelID.String(),
		Description:      db.Description,
		Quantity:         db.Quantity,
		WeightKg:         db.WeightKg,
		LengthCm:         db.LengthCm,
		WidthCm:          db.WidthCm,
		HeightCm:         db.HeightCm,
		VolumetricWeight: db.VolumetricWeight,
		BillableWeight:   db.BillableWeight,
		UnitPrice:        db.UnitPrice,
		ContentType:      db.ContentType,
		Notes:            db.Notes,
		CreatedAt:        db.CreatedAt,
	}
}

// FromDomain convierte domain.ParcelItem a DBParcelItem
func (db *DBParcelItem) FromDomain(item domain.ParcelItem) error {
	id, err := uuid.Parse(item.ID)
	if err != nil && item.ID != "" {
		return err
	}
	if item.ID == "" {
		id = uuid.New()
	}

	parcelID, err := uuid.Parse(item.ParcelID)
	if err != nil {
		return err
	}

	*db = DBParcelItem{
		ID:               id,
		ParcelID:         parcelID,
		Description:      item.Description,
		Quantity:         item.Quantity,
		WeightKg:         item.WeightKg,
		LengthCm:         item.LengthCm,
		WidthCm:          item.WidthCm,
		HeightCm:         item.HeightCm,
		VolumetricWeight: item.VolumetricWeight,
		BillableWeight:   item.BillableWeight,
		UnitPrice:        item.UnitPrice,
		ContentType:      item.ContentType,
		Notes:            item.Notes,
		CreatedAt:        item.CreatedAt,
	}
	return nil
}

// BeforeCreate hook de GORM
func (db *DBParcelItem) BeforeCreate(tx *gorm.DB) error {
	if db.ID == uuid.Nil {
		db.ID = uuid.New()
	}
	return nil
}
