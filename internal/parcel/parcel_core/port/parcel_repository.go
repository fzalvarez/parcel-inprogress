package port

import (
	"context"
	"time"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
)

type ListParcelFilters struct {
	Status              *domain.ParcelStatus
	VehicleID           *string
	OriginOfficeID      *string
	DestinationOfficeID *string
	SenderPersonID      *string
	RecipientPersonID   *string
	FromCreatedAt       *time.Time
	ToCreatedAt         *time.Time
	Limit               int
	Offset              int
	Query               *string // q
}

type ParcelRepository interface {
	Create(ctx context.Context, p domain.Parcel) (uuid.UUID, error)
	GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*domain.Parcel, error)
	UpdateRegistered(ctx context.Context, tenantID string, id uuid.UUID, registeredAtUTC time.Time, userID string, userName string) (*domain.Parcel, error)
	UpdateBoarded(ctx context.Context, tenantID string, id uuid.UUID, boardedAtUTC time.Time, vehicleID string, tripID *string, departureAt *time.Time, boardedByUserID *string) (*domain.Parcel, error)
	ListByFilters(ctx context.Context, tenantID string, f ListParcelFilters) ([]domain.Parcel, error)
	List(ctx context.Context, tenantID string, f ListParcelFilters) (items []domain.Parcel, count int, err error)
	UpdateDelivered(ctx context.Context, tenantID string, id uuid.UUID, deliveredAtUTC time.Time, deliveredByUserID *string) (*domain.Parcel, error)
	UpdateArrivedDestination(ctx context.Context, tenantID string, id uuid.UUID, arrivedAtUTC time.Time, arrivedByUserID *string) (*domain.Parcel, error)
	UpdateInTransit(ctx context.Context, tenantID string, id uuid.UUID, departedAtUTC time.Time, departedByUserID *string, vehicleID *string) (*domain.Parcel, error)
}
