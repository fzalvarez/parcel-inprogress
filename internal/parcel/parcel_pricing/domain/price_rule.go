package domain

import (
	"time"

	coredomain "ms-parcel-core/internal/parcel/parcel_core/domain"
)

type PriceUnit string

const (
	PriceUnitPerKg   PriceUnit = "PER_KG"
	PriceUnitPerItem PriceUnit = "PER_ITEM"
)

const (
	// WildcardOffice representa "cualquier oficina" en reglas de precios
	WildcardOffice = "*"
)

type PriceRule struct {
	ID                  string
	TenantID            string
	ShipmentType        coredomain.ShipmentType
	OriginOfficeID      string
	DestinationOfficeID string
	Unit                PriceUnit
	Price               float64
	Currency            string
	Active              bool
	Priority            int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
