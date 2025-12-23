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

type RegisterParcelInput struct {
	TenantID string
	UserID   string
	UserName string
	ParcelID uuid.UUID
}

type RegisterParcelUseCase struct {
	repo     port.ParcelRepository
	tracking port.TrackingRecorder
}

func NewRegisterParcelUseCase(repo port.ParcelRepository, tracking port.TrackingRecorder) *RegisterParcelUseCase {
	return &RegisterParcelUseCase{repo: repo, tracking: tracking}
}

func (u *RegisterParcelUseCase) Execute(ctx context.Context, in RegisterParcelInput) (*domain.Parcel, error) {
	if strings.TrimSpace(in.TenantID) == "" || strings.TrimSpace(in.UserID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	p, err := u.repo.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}
	if p.Status != domain.ParcelStatusCreated {
		return nil, apperror.New("invalid_state", "transición de estado inválida", map[string]any{"expected": domain.ParcelStatusCreated, "actual": p.Status}, 409)
	}

	registeredAt := time.Now().UTC()
	updated, err := u.repo.UpdateRegistered(ctx, in.TenantID, in.ParcelID, registeredAt, in.UserID, in.UserName)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	if u.tracking != nil {
		if err := u.tracking.RecordEvent(ctx, in.TenantID, port.TrackingEventDTO{
			ParcelID:   in.ParcelID.String(),
			EventType:  port.EventTypeParcelRegistered,
			OccurredAt: registeredAt,
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata:   map[string]any{},
		}); err != nil {
			// TODO: logger
		}
	}

	return updated, nil
}
