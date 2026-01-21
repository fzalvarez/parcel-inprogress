package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	docdomain "ms-parcel-core/internal/parcel/parcel_documents/domain"
)

// DBPrintRecord representa el modelo de base de datos para PrintRecord
type DBPrintRecord struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ParcelID        string    `gorm:"type:varchar(100);not null;index"`
	TenantID        string    `gorm:"type:varchar(100);not null;index"`
	DocumentType    string    `gorm:"type:varchar(50);not null"`
	PrintedAt       time.Time `gorm:"not null"`
	PrintedByUserID *string   `gorm:"type:varchar(100)"`
}

func (DBPrintRecord) TableName() string {
	return "print_records"
}

// ToDomain convierte DBPrintRecord a docdomain.PrintRecord
func (db *DBPrintRecord) ToDomain() docdomain.PrintRecord {
	return docdomain.PrintRecord{
		ID:              db.ID.String(),
		ParcelID:        db.ParcelID,
		DocumentType:    docdomain.DocumentType(db.DocumentType),
		PrintedAt:       db.PrintedAt,
		PrintedByUserID: db.PrintedByUserID,
	}
}

// FromDomain convierte docdomain.PrintRecord a DBPrintRecord
func (db *DBPrintRecord) FromDomain(rec docdomain.PrintRecord) error {
	id, err := uuid.Parse(rec.ID)
	if err != nil && rec.ID != "" {
		return err
	}
	if rec.ID == "" {
		id = uuid.New()
	}

	*db = DBPrintRecord{
		ID:              id,
		ParcelID:        rec.ParcelID,
		DocumentType:    string(rec.DocumentType),
		PrintedAt:       rec.PrintedAt,
		PrintedByUserID: rec.PrintedByUserID,
	}
	return nil
}

// BeforeCreate hook de GORM
func (db *DBPrintRecord) BeforeCreate(tx *gorm.DB) error {
	if db.ID == uuid.Nil {
		db.ID = uuid.New()
	}
	return nil
}
