package domain

import "time"

type PaymentType string

type Currency string

type PaymentStatus string

type PaymentChannel string

const (
	PaymentTypeCash              PaymentType = "CASH"
	PaymentTypeFOB               PaymentType = "FOB"
	PaymentTypeCard              PaymentType = "CARD"
	PaymentTypeTransfer          PaymentType = "TRANSFER"
	PaymentTypeEWallet           PaymentType = "EWALLET"
	PaymentTypeFree              PaymentType = "FREE"
	PaymentTypeCollectOnDelivery PaymentType = "COLLECT_ON_DELIVERY"
)

const (
	CurrencyPEN Currency = "PEN"
	CurrencyUSD Currency = "USD"
)

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
)

const (
	PaymentChannelCounter PaymentChannel = "COUNTER"
	PaymentChannelWeb     PaymentChannel = "WEB"
)

type ParcelPayment struct {
	ID           string
	TenantID     string
	ParcelID     string
	PaymentType  PaymentType
	Currency     Currency
	Amount       float64
	Notes        *string
	Status       PaymentStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PaidAt       *time.Time
	PaidByUserID *string

	Channel      PaymentChannel
	OfficeID     *string
	CashboxID    *string
	SellerUserID *string
}
