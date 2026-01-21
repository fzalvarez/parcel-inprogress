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
	LengthCm    *float64
	WidthCm     *float64
	HeightCm    *float64
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

func (u *AddParcelItemUseCase) Execute(ctx context.Context, in AddParcelItemInput) (*domain.ParcelItem, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	parcel, err := u.parcelReader.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if parcel == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	defaults := coreport.ParcelOptions{
		RequirePackageKey:       true,
		UsePriceTable:           true,
		UseVolumetricWeight:     false,
		VolumetricDivisor:       6000,
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

	// Cálculo de peso volumétrico y facturable
	var volumetricWeight *float64
	billableWeight := in.WeightKg

	if opts.UseVolumetricWeight && in.LengthCm != nil && in.WidthCm != nil && in.HeightCm != nil {
		divisor := float64(opts.VolumetricDivisor)
		if divisor <= 0 {
			divisor = 6000
		}
		vw := (*in.LengthCm * *in.WidthCm * *in.HeightCm) / divisor
		volumetricWeight = &vw

		if vw > in.WeightKg {
			billableWeight = vw
		}
	}

	unitPrice := in.UnitPrice

	if opts.UsePriceTable {
		if u.priceRules == nil {
			return nil, apperror.New("price_rule_not_found", "regla de precios no configurada", nil, 409)
		}

		rule, err := u.priceRules.FindMatch(ctx, in.TenantID, string(parcel.ShipmentType), parcel.OriginOfficeID, parcel.DestinationOfficeID)
		if err != nil {
			return nil, err
		}
		if rule == nil {
			if !opts.AllowManualPrice || in.UnitPrice <= 0 {
				return nil, apperror.New("price_rule_not_found", "regla de precios no encontrada para esta ruta. Defina una regla específica o use comodín (*)", map[string]any{
					"shipment_type":         parcel.ShipmentType,
					"origin_office_id":      parcel.OriginOfficeID,
					"destination_office_id": parcel.DestinationOfficeID,
					"hint":                  "Puede crear una regla global usando '*' como origin_office_id o destination_office_id",
				}, 409)
			}
			// Si permite precio manual y el usuario lo envió, continuamos sin regla
		} else {
			suggested := 0.0
			switch rule.Unit {
			case pricingdomain.PriceUnitPerItem:
				suggested = rule.Price * float64(in.Quantity)
			case pricingdomain.PriceUnitPerKg:
				suggested = rule.Price * billableWeight
			default:
				return nil, apperror.New("validation_error", "unit inválido", map[string]any{"unit": rule.Unit}, 400)
			}

			if !opts.AllowOverridePriceTable {
				if !opts.AllowManualPrice && in.UnitPrice > 0 {
					return nil, apperror.New("manual_price_disabled", "precio manual deshabilitado", nil, 409)
				}
				unitPrice = suggested
			} else {
				if in.UnitPrice <= 0 {
					unitPrice = suggested
				}
			}
		}
	}

	in.UnitPrice = unitPrice

	item := domain.ParcelItem{
		ID:               uuid.NewString(),
		ParcelID:         in.ParcelID.String(),
		Description:      in.Description,
		Quantity:         in.Quantity,
		WeightKg:         in.WeightKg,
		LengthCm:         in.LengthCm,
		WidthCm:          in.WidthCm,
		HeightCm:         in.HeightCm,
		VolumetricWeight: volumetricWeight,
		BillableWeight:   billableWeight,
		UnitPrice:        in.UnitPrice,
		ContentType:      in.ContentType,
		Notes:            in.Notes,
		CreatedAt:        time.Now().UTC(),
	}

	id, err := u.repo.Add(ctx, in.TenantID, item)
	if err != nil {
		return nil, err
	}

	item.ID = id.String()

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

	return &item, nil
}
