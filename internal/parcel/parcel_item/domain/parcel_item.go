package domain

import "time"

type ParcelItem struct {
	ID          string
	ParcelID    string
	Description string
	Quantity    int
	WeightKg    float64
	UnitPrice   float64
	ContentType *string
	Notes       *string
	CreatedAt   time.Time
}
