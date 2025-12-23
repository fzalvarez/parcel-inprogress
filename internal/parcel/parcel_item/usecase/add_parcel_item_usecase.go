package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	coredomain "ms-parcel-core/internal/parcel/parcel_core/domain"
	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_item/domain"
	"ms-parcel-core/internal/parcel/parcel_item/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type AddParcelItemInput struct {
	TenantID    string
	UserID      string
	UserName    string
	ParcelID    uuid.UUID
	Description string
	Quantity    int
	WeightKg    float64
	UnitPrice   float64
	ContentType *string
	Notes       *string
}

type AddParcelItemUseCase struct {
	parcelReader coreport.ParcelReader
	repo         port.ParcelItemRepository
	tracking     coreport.TrackingRecorder
}

func NewAddParcelItemUseCase(parcelReader coreport.ParcelReader, repo port.ParcelItemRepository, tracking coreport.TrackingRecorder) *AddParcelItemUseCase {
	return &AddParcelItemUseCase{parcelReader: parcelReader, repo: repo, tracking: tracking}
}

func (u *AddParcelItemUseCase) Execute(ctx context.Context, in AddParcelItemInput) (uuid.UUID, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return uuid.Nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return uuid.Nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	p, err := u.parcelReader.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return uuid.Nil, err
	}
	if p == nil {
		return uuid.Nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	allowed := p.Status == coredomain.ParcelStatusCreated || p.Status == coredomain.ParcelStatusRegistered
	if !allowed {
		return uuid.Nil, apperror.New("invalid_state", "no se pueden modificar items en este estado", map[string]any{"allowed": []coredomain.ParcelStatus{coredomain.ParcelStatusCreated, coredomain.ParcelStatusRegistered}, "actual": p.Status}, 409)
	}

	// TODO: validar desde TENANT-CONFIG si el precio manual está permitido.

	item := domain.ParcelItem{
		ID:          uuid.NewString(),
		ParcelID:    in.ParcelID.String(),
		Description: in.Description,
		Quantity:    in.Quantity,
		WeightKg:    in.WeightKg,
		UnitPrice:   in.UnitPrice,
		ContentType: in.ContentType,
		Notes:       in.Notes,
		CreatedAt:   time.Now().UTC(),
	}

	id, err := u.repo.Add(ctx, in.TenantID, item)
	if err != nil {
		return uuid.Nil, err
	}

	if u.tracking != nil {
		_ = u.tracking.RecordEvent(ctx, in.TenantID, coreport.TrackingEventDTO{
			ParcelID:   in.ParcelID.String(),
			EventType:  "PARCEL_ITEM_ADDED",
			OccurredAt: time.Now().UTC(),
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata: map[string]any{
				"item_id":   id.String(),
				"quantity":  in.Quantity,
				"weight_kg": in.WeightKg,
			},
		})
		// TODO: logger si falla
	}

	return id, nil
}
