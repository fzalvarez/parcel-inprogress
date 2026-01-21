package domain

import "time"

type ParcelItem struct {
	ID               string
	ParcelID         string
	Description      string
	Quantity         int
	WeightKg         float64
	LengthCm         *float64
	WidthCm          *float64
	HeightCm         *float64
	VolumetricWeight *float64
	BillableWeight   float64
	UnitPrice        float64
	ContentType      *string
	Notes            *string
	CreatedAt        time.Time
}
