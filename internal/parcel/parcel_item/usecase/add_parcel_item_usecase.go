package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_item/domain"
	"ms-parcel-core/internal/parcel/parcel_item/port"
	pricingdomain "ms-parcel-core/internal/parcel/parcel_pricing/domain"
	pricingport "ms-parcel-core/internal/parcel/parcel_pricing/port"
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
	priceRules      pricingport.PriceRuleRepository
}

func NewAddParcelItemUseCase(parcelReader coreport.ParcelReader, repo port.ParcelItemRepository, tracking coreport.TrackingRecorder, optionsProvider coreport.TenantOptionsProvider, priceRules pricingport.PriceRuleRepository) *AddParcelItemUseCase {
	return &AddParcelItemUseCase{parcelReader: parcelReader, repo: repo, tracking: tracking, optionsProvider: optionsProvider, priceRules: priceRules}
}

func (u *AddParcelItemUseCase) Execute(ctx context.Context, in AddParcelItemInput) (uuid.UUID, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return uuid.Nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return uuid.Nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	parcel, err := u.parcelReader.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return uuid.Nil, err
	}
	if parcel == nil {
		return uuid.Nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	defaults := coreport.ParcelOptions{
		RequirePackageKey:       true,
		UsePriceTable:           true,
		AllowManualPrice:        false,
		AllowOverridePriceTable: true,
		AllowPayInDestination:   false,
		MaxPrints:               1,
		AllowReprint:            false,
		ReprintFeeEnabled:       false,
	}
	opts := defaults
	if u.optionsProvider != nil {
		if o, err := u.optionsProvider.GetParcelOptions(ctx, in.TenantID); err == nil {
			opts = o
		} else {
			// TODO: logger
		}
	}

	unitPrice := in.UnitPrice

	if opts.UsePriceTable {
		if u.priceRules == nil {
			return uuid.Nil, apperror.New("price_rule_not_found", "regla de precios no configurada", nil, 409)
		}

		rule, err := u.priceRules.FindMatch(ctx, in.TenantID, string(parcel.ShipmentType), parcel.OriginOfficeID, parcel.DestinationOfficeID)
		if err != nil {
			return uuid.Nil, err
		}
		if rule == nil {
			return uuid.Nil, apperror.New("price_rule_not_found", "regla de precios no encontrada", map[string]any{"shipment_type": parcel.ShipmentType, "origin_office_id": parcel.OriginOfficeID, "destination_office_id": parcel.DestinationOfficeID}, 409)
		}

		suggested := 0.0
		switch rule.Unit {
		case pricingdomain.PriceUnitPerItem:
			suggested = rule.Price * float64(in.Quantity)
		case pricingdomain.PriceUnitPerKg:
			suggested = rule.Price * in.WeightKg
		default:
			return uuid.Nil, apperror.New("validation_error", "unit inválido", map[string]any{"unit": rule.Unit}, 400)
		}

		if !opts.AllowOverridePriceTable {
			if !opts.AllowManualPrice && in.UnitPrice > 0 {
				return uuid.Nil, apperror.New("manual_price_disabled", "precio manual deshabilitado", nil, 409)
			}
			unitPrice = suggested
		} else {
			if in.UnitPrice <= 0 {
				unitPrice = suggested
			}
		}
	}

	in.UnitPrice = unitPrice

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
