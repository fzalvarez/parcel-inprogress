package dto

type CreateParcelRequest struct {
	ShipmentType        string  `json:"shipment_type" binding:"required"`
	OriginOfficeID      string  `json:"origin_office_id" binding:"required,uuid"`
	DestinationOfficeID string  `json:"destination_office_id" binding:"required,uuid"`
	SenderPersonID      string  `json:"sender_person_id" binding:"required,uuid"`
	RecipientPersonID   string  `json:"recipient_person_id" binding:"required,uuid"`
	Notes               *string `json:"notes" binding:"omitempty,max=200"`
	PackageKey          string  `json:"package_key" binding:"omitempty,max=50"`
	PackageKeyConfirm   string  `json:"package_key_confirm" binding:"omitempty,max=50"`
}

type CreateParcelResponse struct {
	ID                  string  `json:"id"`
	Status              string  `json:"status"`
	ShipmentType        string  `json:"shipment_type"`
	OriginOfficeID      string  `json:"origin_office_id"`
	DestinationOfficeID string  `json:"destination_office_id"`
	SenderPersonID      string  `json:"sender_person_id"`
	RecipientPersonID   string  `json:"recipient_person_id"`
	Notes               *string `json:"notes,omitempty"`
	CreatedAt           string  `json:"created_at"`
	RegisteredAt        *string `json:"registered_at,omitempty"`
	BoardedVehicleID    *string `json:"boarded_vehicle_id,omitempty"`
	BoardedTripID       *string `json:"boarded_trip_id,omitempty"`
	BoardedDepartureAt  *string `json:"boarded_departure_at,omitempty"`
	BoardedAt           *string `json:"boarded_at,omitempty"`
	BoardedByUserID     *string `json:"boarded_by_user_id,omitempty"`
	DepartedAt          *string `json:"departed_at,omitempty"`
	DepartedByUserID    *string `json:"departed_by_user_id,omitempty"`
	ArrivedAt           *string `json:"arrived_at,omitempty"`
	ArrivedByUserID     *string `json:"arrived_by_user_id,omitempty"`
	DeliveredAt         *string `json:"delivered_at,omitempty"`
	DeliveredByUserID   *string `json:"delivered_by_user_id,omitempty"`
}

type ParcelListPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type ParcelListResponse struct {
	Items      []CreateParcelResponse `json:"items"`
	Pagination ParcelListPagination   `json:"pagination"`
}

type BoardParcelRequest struct {
	VehicleID      string  `json:"vehicle_id" binding:"required,uuid"`
	TripID         *string `json:"trip_id" binding:"omitempty,uuid"`
	DepartureAt    *string `json:"departure_at" binding:"omitempty"`
	OriginOfficeID *string `json:"origin_office_id" binding:"omitempty,uuid"`
}

type DepartParcelRequest struct {
	DepartureOfficeID string  `json:"departure_office_id" binding:"omitempty,uuid"`
	VehicleID         *string `json:"vehicle_id" binding:"omitempty,uuid"`
	DepartedAt        *string `json:"departed_at" binding:"omitempty"`
}

type ArriveParcelRequest struct {
	DestinationOfficeID string `json:"destination_office_id" binding:"omitempty,uuid"`
}

type DeliverParcelRequest struct {
	PackageKey string `json:"package_key" binding:"omitempty,max=50"`
}
