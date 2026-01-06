package domain

import "time"

type ParcelStatus string

const (
	ParcelStatusCreated            ParcelStatus = "CREATED"
	ParcelStatusRegistered         ParcelStatus = "REGISTERED"
	ParcelStatusBoarded            ParcelStatus = "EMBARCADO"
	ParcelStatusInTransit          ParcelStatus = "EN_TRANSITO"
	ParcelStatusArrivedDestination ParcelStatus = "EN_OFICINA_DESTINO"
	ParcelStatusDelivered          ParcelStatus = "ENTREGADO"
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
	RegisteredAt         *time.Time

	BoardedVehicleID   *string
	BoardedTripID      *string
	BoardedDepartureAt *time.Time
	BoardedAt          *time.Time
	BoardedByUserID    *string
	DeliveredAt        *time.Time
	DeliveredByUserID  *string
	ArrivedAt          *time.Time
	ArrivedByUserID    *string
	DepartedAt         *time.Time
	DepartedByUserID   *string
	TrackingCode       string
}
