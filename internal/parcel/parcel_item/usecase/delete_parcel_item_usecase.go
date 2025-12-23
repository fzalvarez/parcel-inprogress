package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	coredomain "ms-parcel-core/internal/parcel/parcel_core/domain"
	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_item/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type DeleteParcelItemInput struct {
	TenantID string
	UserID   string
	UserName string
	ParcelID uuid.UUID
	ItemID   uuid.UUID
}

type DeleteParcelItemUseCase struct {
	parcelReader coreport.ParcelReader
	repo         port.ParcelItemRepository
	tracking     coreport.TrackingRecorder
}

func NewDeleteParcelItemUseCase(parcelReader coreport.ParcelReader, repo port.ParcelItemRepository, tracking coreport.TrackingRecorder) *DeleteParcelItemUseCase {
	return &DeleteParcelItemUseCase{parcelReader: parcelReader, repo: repo, tracking: tracking}
}

func (u *DeleteParcelItemUseCase) Execute(ctx context.Context, in DeleteParcelItemInput) error {
	if strings.TrimSpace(in.TenantID) == "" {
		return apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}
	if in.ItemID == uuid.Nil {
		return apperror.NewBadRequest("validation_error", "item_id inválido", map[string]any{"field": "item_id"})
	}

	p, err := u.parcelReader.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return err
	}
	if p == nil {
		return apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	allowed := p.Status == coredomain.ParcelStatusCreated || p.Status == coredomain.ParcelStatusRegistered
	if !allowed {
		return apperror.New("invalid_state", "no se pueden modificar items en este estado", map[string]any{"allowed": []coredomain.ParcelStatus{coredomain.ParcelStatusCreated, coredomain.ParcelStatusRegistered}, "actual": p.Status}, 409)
	}

	if err := u.repo.Delete(ctx, in.TenantID, in.ParcelID, in.ItemID); err != nil {
		return err
	}

	if u.tracking != nil {
		_ = u.tracking.RecordEvent(ctx, in.TenantID, coreport.TrackingEventDTO{
			ParcelID:   in.ParcelID.String(),
			EventType:  "PARCEL_ITEM_REMOVED",
			OccurredAt: time.Now().UTC(),
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata:   map[string]any{"item_id": in.ItemID.String()},
		})
		// TODO: logger si falla
	}

	return nil
}
