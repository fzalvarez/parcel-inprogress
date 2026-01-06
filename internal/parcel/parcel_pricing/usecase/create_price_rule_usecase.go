package usecase

import (
	"context"
	"strings"

	coredomain "ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_pricing/domain"
	"ms-parcel-core/internal/parcel/parcel_pricing/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type CreatePriceRuleUseCase struct {
	repo port.PriceRuleRepository
}

func NewCreatePriceRuleUseCase(repo port.PriceRuleRepository) *CreatePriceRuleUseCase {
	return &CreatePriceRuleUseCase{repo: repo}
}

type CreatePriceRuleInput struct {
	TenantID            string
	ShipmentType        string
	OriginOfficeID      string
	DestinationOfficeID string
	Unit                string
	Price               float64
	Currency            string
	Active              bool
}

func (u *CreatePriceRuleUseCase) Execute(ctx context.Context, in CreatePriceRuleInput) (*domain.PriceRule, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.Price <= 0 {
		return nil, apperror.NewBadRequest("validation_error", "price inválido", map[string]any{"field": "price"})
	}
	switch strings.TrimSpace(in.Unit) {
	case string(domain.PriceUnitPerKg), string(domain.PriceUnitPerItem):
	default:
		return nil, apperror.NewBadRequest("validation_error", "unit inválido", map[string]any{"field": "unit"})
	}
	switch strings.TrimSpace(in.Currency) {
	case "PEN", "USD":
	default:
		return nil, apperror.NewBadRequest("validation_error", "currency inválido", map[string]any{"field": "currency"})
	}

	r := domain.PriceRule{
		ShipmentType:        coredomain.ShipmentType(in.ShipmentType),
		OriginOfficeID:      strings.TrimSpace(in.OriginOfficeID),
		DestinationOfficeID: strings.TrimSpace(in.DestinationOfficeID),
		Unit:                domain.PriceUnit(strings.TrimSpace(in.Unit)),
		Price:               in.Price,
		Currency:            strings.TrimSpace(in.Currency),
		Active:              in.Active,
	}

	return u.repo.Create(ctx, in.TenantID, r)
}
