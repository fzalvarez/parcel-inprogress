package domain

import "time"

type ParcelStatus string

const (
	ParcelStatusCreated ParcelStatus = "CREATED"
)

type ShipmentType string

const (
	ShipmentTypeBus      ShipmentType = "BUS"
	ShipmentTypeCarguero ShipmentType = "CARGUERO"
)

type Parcel struct {
	ID                   string
	TenantID             string
	OriginOfficeID       string
	DestinationOfficeID  string
	SenderPersonID       string
	RecipientPersonID    string
	ShipmentType         ShipmentType
	Notes                *string
	PackageKeyHashSHA256 string
	Status               ParcelStatus
	CreatedByUserID      string
	CreatedByUserName    string
	CreatedAt            time.Time
}
