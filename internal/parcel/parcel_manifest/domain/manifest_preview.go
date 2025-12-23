package domain

type ParcelSummary struct {
	ParcelID          string
	Status            string
	SenderPersonID    string
	RecipientPersonID string
	Notes             *string
}

type ManifestTotals struct {
	CountParcels int
}

type ManifestPreview struct {
	VehicleID           string
	OriginOfficeID      string
	DestinationOfficeID string
	Parcels             []ParcelSummary
	Totals              ManifestTotals
}
