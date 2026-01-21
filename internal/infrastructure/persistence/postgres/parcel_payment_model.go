package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	paymentdomain "ms-parcel-core/internal/parcel/parcel_payment/domain"
)

// DBParcelPayment representa el modelo de base de datos para ParcelPayment
type DBParcelPayment struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ParcelID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	TenantID     string    `gorm:"type:varchar(100);not null;index"`
	PaymentType  string    `gorm:"type:varchar(50);not null"`
	Status       string    `gorm:"type:varchar(50);not null"`
	Amount       float64   `gorm:"type:decimal(10,2);not null"`
	Currency     string    `gorm:"type:varchar(3);not null"`
	Channel      string    `gorm:"type:varchar(50)"`
	OfficeID     *string   `gorm:"type:varchar(100)"`
	CashboxID    *string   `gorm:"type:varchar(100)"`
	SellerUserID *string   `gorm:"type:varchar(100)"`
	Notes        *string   `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	PaidAt       *time.Time
	PaidByUserID *string `gorm:"type:varchar(100)"`
}

func (DBParcelPayment) TableName() string {
	return "parcel_payments"
}

// ToDomain convierte DBParcelPayment a paymentdomain.ParcelPayment
func (db *DBParcelPayment) ToDomain() paymentdomain.ParcelPayment {
	return paymentdomain.ParcelPayment{
		ID:           db.ID.String(),
		ParcelID:     db.ParcelID.String(),
		TenantID:     db.TenantID,
		PaymentType:  paymentdomain.PaymentType(db.PaymentType),
		Status:       paymentdomain.PaymentStatus(db.Status),
		Amount:       db.Amount,
		Currency:     paymentdomain.Currency(db.Currency),
		Channel:      paymentdomain.PaymentChannel(db.Channel),
		OfficeID:     db.OfficeID,
		CashboxID:    db.CashboxID,
		SellerUserID: db.SellerUserID,
		Notes:        db.Notes,
		CreatedAt:    db.CreatedAt,
		UpdatedAt:    db.UpdatedAt,
		PaidAt:       db.PaidAt,
		PaidByUserID: db.PaidByUserID,
	}
}

// FromDomain convierte paymentdomain.ParcelPayment a DBParcelPayment
func (db *DBParcelPayment) FromDomain(p paymentdomain.ParcelPayment) error {
	id, err := uuid.Parse(p.ID)
	if err != nil && p.ID != "" {
		return err
	}
	if p.ID == "" {
		id = uuid.New()
	}

	parcelID, err := uuid.Parse(p.ParcelID)
	if err != nil {
		return err
	}

	*db = DBParcelPayment{
		ID:           id,
		ParcelID:     parcelID,
		TenantID:     p.TenantID,
		PaymentType:  string(p.PaymentType),
		Status:       string(p.Status),
		Amount:       p.Amount,
		Currency:     string(p.Currency),
		Channel:      string(p.Channel),
		OfficeID:     p.OfficeID,
		CashboxID:    p.CashboxID,
		SellerUserID: p.SellerUserID,
		Notes:        p.Notes,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		PaidAt:       p.PaidAt,
		PaidByUserID: p.PaidByUserID,
	}
	return nil
}

// BeforeCreate hook de GORM
func (db *DBParcelPayment) BeforeCreate(tx *gorm.DB) error {
	if db.ID == uuid.Nil {
		db.ID = uuid.New()
	}
	return nil
}
