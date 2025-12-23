package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type BoardParcelInput struct {
	TenantID    string
	UserID      string
	UserName    string
	ParcelID    uuid.UUID
	VehicleID   uuid.UUID
	TripID      *uuid.UUID
	DepartureAt *time.Time
}

type BoardParcelUseCase struct {
	repo     port.ParcelRepository
	tracking port.TrackingRecorder
}

func NewBoardParcelUseCase(repo port.ParcelRepository, tracking port.TrackingRecorder) *BoardParcelUseCase {
	return &BoardParcelUseCase{repo: repo, tracking: tracking}
}

func (u *BoardParcelUseCase) Execute(ctx context.Context, in BoardParcelInput) (*domain.Parcel, error) {
	if strings.TrimSpace(in.TenantID) == "" || strings.TrimSpace(in.UserID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}
	if in.VehicleID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "vehicle_id inválido", map[string]any{"field": "vehicle_id"})
	}

	p, err := u.repo.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}
	if p.Status != domain.ParcelStatusRegistered {
		return nil, apperror.New("invalid_state", "transición de estado inválida", map[string]any{"expected": domain.ParcelStatusRegistered, "actual": p.Status}, 409)
	}

	boardedAt := time.Now().UTC()
	vehicleIDStr := in.VehicleID.String()

	var tripIDStr *string
	if in.TripID != nil && *in.TripID != uuid.Nil {
		s := in.TripID.String()
		tripIDStr = &s
	}

	boardedBy := strings.TrimSpace(in.UserID)
	boardedByPtr := &boardedBy

	updated, err := u.repo.UpdateBoarded(ctx, in.TenantID, in.ParcelID, boardedAt, vehicleIDStr, tripIDStr, in.DepartureAt, boardedByPtr)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	if u.tracking != nil {
		md := map[string]any{"vehicle_id": vehicleIDStr}
		if tripIDStr != nil {
			md["trip_id"] = *tripIDStr
		}
		if in.DepartureAt != nil {
			md["departure_at"] = in.DepartureAt.UTC().Format(time.RFC3339)
		}
		if err := u.tracking.RecordEvent(ctx, in.TenantID, port.TrackingEventDTO{
			ParcelID:   in.ParcelID.String(),
			EventType:  port.EventTypeParcelBoarded,
			OccurredAt: boardedAt,
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata:   md,
		}); err != nil {
			// TODO: logger
		}
	}

	return updated, nil
}
