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

type DepartParcelInput struct {
	TenantID          string
	UserID            string
	UserName          string
	ParcelID          uuid.UUID
	DepartureOfficeID string
	VehicleID         *uuid.UUID
	DepartedAt        *time.Time
}

type DepartParcelUseCase struct {
	repo     port.ParcelRepository
	tracking port.TrackingRecorder
}

func NewDepartParcelUseCase(repo port.ParcelRepository, tracking port.TrackingRecorder) *DepartParcelUseCase {
	return &DepartParcelUseCase{repo: repo, tracking: tracking}
}

func (u *DepartParcelUseCase) Execute(ctx context.Context, in DepartParcelInput) (*domain.Parcel, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if strings.TrimSpace(in.UserID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}
	if strings.TrimSpace(in.DepartureOfficeID) == "" {
		return nil, apperror.NewBadRequest("validation_error", "departure_office_id requerido", map[string]any{"field": "departure_office_id"})
	}

	p, err := u.repo.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	if p.OriginOfficeID != in.DepartureOfficeID {
		return nil, apperror.New("origin_mismatch", "departure_office_id no coincide", map[string]any{"expected": p.OriginOfficeID, "actual": in.DepartureOfficeID}, 409)
	}

	if p.Status != domain.ParcelStatusBoarded {
		return nil, apperror.New("invalid_state", "transición de estado inválida", map[string]any{"expected": domain.ParcelStatusBoarded, "actual": p.Status}, 409)
	}

	var vehicleIDStr *string
	if in.VehicleID != nil && *in.VehicleID != uuid.Nil {
		v := in.VehicleID.String()
		if p.BoardedVehicleID != nil && *p.BoardedVehicleID != v {
			return nil, apperror.New("vehicle_mismatch", "vehicle_id no coincide", map[string]any{"expected": *p.BoardedVehicleID, "actual": v}, 409)
		}
		vehicleIDStr = &v
		// Si está vacío en parcel, se setea en UpdateInTransit (MVP)
		if p.BoardedVehicleID == nil {
			// ok
		}
	}

	departedAt := time.Now().UTC()
	if in.DepartedAt != nil {
		departedAt = in.DepartedAt.UTC()
	}
	by := strings.TrimSpace(in.UserID)
	byPtr := &by

	updated, err := u.repo.UpdateInTransit(ctx, in.TenantID, in.ParcelID, departedAt, byPtr, vehicleIDStr)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	if u.tracking != nil {
		md := map[string]any{
			"departure_office_id": in.DepartureOfficeID,
			"departed_at":         departedAt.UTC().Format(time.RFC3339),
			"departed_by_user_id": by,
		}
		if vehicleIDStr != nil {
			md["vehicle_id"] = *vehicleIDStr
		}

		if err := u.tracking.RecordEvent(ctx, in.TenantID, port.TrackingEventDTO{
			ParcelID:   in.ParcelID.String(),
			EventType:  port.EventTypeParcelInTransit,
			OccurredAt: departedAt,
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata:   md,
		}); err != nil {
			// TODO: logger
		}
	}

	return updated, nil
}
