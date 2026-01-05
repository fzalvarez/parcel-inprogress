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
	parcelReader    coreport.ParcelReader
	repo            port.ParcelItemRepository
	tracking        coreport.TrackingRecorder
	optionsProvider coreport.TenantOptionsProvider
}

func NewAddParcelItemUseCase(parcelReader coreport.ParcelReader, repo port.ParcelItemRepository, tracking coreport.TrackingRecorder, optionsProvider coreport.TenantOptionsProvider) *AddParcelItemUseCase {
	return &AddParcelItemUseCase{parcelReader: parcelReader, repo: repo, tracking: tracking, optionsProvider: optionsProvider}
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

	defaults := coreport.ParcelOptions{
		RequirePackageKey:       true,
		UsePriceTable:           true,
		AllowManualPrice:        false,
		AllowOverridePriceTable: true,
		AllowPayInDestination:   false,
	}
	opts := defaults
	if u.optionsProvider != nil {
		if o, err := u.optionsProvider.GetParcelOptions(ctx, in.TenantID); err == nil {
			opts = o
		} else {
			// TODO: logger
		}
	}

	if opts.UsePriceTable && !opts.AllowManualPrice {
		if in.UnitPrice > 0 {
			return uuid.Nil, apperror.New("manual_price_disabled", "precio manual deshabilitado", nil, 409)
		}
		// TODO: cuando exista tabla real, calcular unit_price automáticamente; override permitido si opts.AllowOverridePriceTable
	}

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
