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

type ArriveParcelInput struct {
	TenantID            string
	UserID              string
	UserName            string
	ParcelID            uuid.UUID
	DestinationOfficeID string
}

type ArriveParcelUseCase struct {
	repo     port.ParcelRepository
	tracking port.TrackingRecorder
}

func NewArriveParcelUseCase(repo port.ParcelRepository, tracking port.TrackingRecorder) *ArriveParcelUseCase {
	return &ArriveParcelUseCase{repo: repo, tracking: tracking}
}

func (u *ArriveParcelUseCase) Execute(ctx context.Context, in ArriveParcelInput) (*domain.Parcel, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if strings.TrimSpace(in.UserID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}
	if strings.TrimSpace(in.DestinationOfficeID) == "" {
		return nil, apperror.NewBadRequest("validation_error", "destination_office_id requerido", map[string]any{"field": "destination_office_id"})
	}

	p, err := u.repo.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	if p.DestinationOfficeID != in.DestinationOfficeID {
		return nil, apperror.New("destination_mismatch", "destination_office_id no coincide", map[string]any{"expected": p.DestinationOfficeID, "actual": in.DestinationOfficeID}, 409)
	}

	// MVP: permitir desde EMBARCADO (EN_TRANSITO no existe aún)
	if p.Status != domain.ParcelStatusBoarded {
		return nil, apperror.New("invalid_state", "transición de estado inválida", map[string]any{"expected": domain.ParcelStatusBoarded, "actual": p.Status}, 409)
	}

	arrivedAt := time.Now().UTC()
	by := strings.TrimSpace(in.UserID)
	byPtr := &by

	updated, err := u.repo.UpdateArrivedDestination(ctx, in.TenantID, in.ParcelID, arrivedAt, byPtr)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	if u.tracking != nil {
		if err := u.tracking.RecordEvent(ctx, in.TenantID, port.TrackingEventDTO{
			ParcelID:   in.ParcelID.String(),
			EventType:  port.EventTypeParcelArrivedDestination,
			OccurredAt: arrivedAt,
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata: map[string]any{
				"destination_office_id": in.DestinationOfficeID,
				"arrived_at":            arrivedAt.UTC().Format(time.RFC3339),
				"arrived_by_user_id":    by,
			},
		}); err != nil {
			// TODO: logger
		}
	}

	return updated, nil
}
