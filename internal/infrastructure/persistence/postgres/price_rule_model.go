package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	coredomain "ms-parcel-core/internal/parcel/parcel_core/domain"
	pricingdomain "ms-parcel-core/internal/parcel/parcel_pricing/domain"
)

// DBPriceRule representa el modelo de base de datos para PriceRule
type DBPriceRule struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TenantID            string    `gorm:"type:varchar(100);not null;index"`
	ShipmentType        string    `gorm:"type:varchar(50);not null"`
	OriginOfficeID      string    `gorm:"type:varchar(100);not null;index"`
	DestinationOfficeID string    `gorm:"type:varchar(100);not null;index"`
	Unit                string    `gorm:"type:varchar(50);not null"`
	Price               float64   `gorm:"type:decimal(10,2);not null"`
	Currency            string    `gorm:"type:varchar(3);not null"`
	Priority            int       `gorm:"not null;default:0"`
	Active              bool      `gorm:"not null;default:true"`
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
}

func (DBPriceRule) TableName() string {
	return "price_rules"
}

// ToDomain convierte DBPriceRule a pricingdomain.PriceRule
func (db *DBPriceRule) ToDomain() pricingdomain.PriceRule {
	return pricingdomain.PriceRule{
		ID:                  db.ID.String(),
		TenantID:            db.TenantID,
		ShipmentType:        coredomain.ShipmentType(db.ShipmentType),
		OriginOfficeID:      db.OriginOfficeID,
		DestinationOfficeID: db.DestinationOfficeID,
		Unit:                pricingdomain.PriceUnit(db.Unit),
		Price:               db.Price,
		Currency:            db.Currency,
		Priority:            db.Priority,
		Active:              db.Active,
		CreatedAt:           db.CreatedAt,
		UpdatedAt:           db.UpdatedAt,
	}
}

// FromDomain convierte pricingdomain.PriceRule a DBPriceRule
func (db *DBPriceRule) FromDomain(rule pricingdomain.PriceRule) error {
	id, err := uuid.Parse(rule.ID)
	if err != nil && rule.ID != "" {
		return err
	}
	if rule.ID == "" {
		id = uuid.New()
	}

	*db = DBPriceRule{
		ID:                  id,
		TenantID:            rule.TenantID,
		ShipmentType:        string(rule.ShipmentType),
		OriginOfficeID:      rule.OriginOfficeID,
		DestinationOfficeID: rule.DestinationOfficeID,
		Unit:                string(rule.Unit),
		Price:               rule.Price,
		Currency:            rule.Currency,
		Priority:            rule.Priority,
		Active:              rule.Active,
		CreatedAt:           rule.CreatedAt,
		UpdatedAt:           rule.UpdatedAt,
	}
	return nil
}

// BeforeCreate hook de GORM
func (db *DBPriceRule) BeforeCreate(tx *gorm.DB) error {
	if db.ID == uuid.Nil {
		db.ID = uuid.New()
	}
	return nil
}
