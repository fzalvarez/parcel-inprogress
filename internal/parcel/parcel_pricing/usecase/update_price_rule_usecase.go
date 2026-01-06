package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	coredomain "ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_pricing/domain"
	"ms-parcel-core/internal/parcel/parcel_pricing/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type UpdatePriceRuleUseCase struct {
	repo port.PriceRuleRepository
}

func NewUpdatePriceRuleUseCase(repo port.PriceRuleRepository) *UpdatePriceRuleUseCase {
	return &UpdatePriceRuleUseCase{repo: repo}
}

type UpdatePriceRuleInput struct {
	TenantID            string
	ID                  uuid.UUID
	ShipmentType        string
	OriginOfficeID      string
	DestinationOfficeID string
	Unit                string
	Price               float64
	Currency            string
	Active              bool
}

func (u *UpdatePriceRuleUseCase) Execute(ctx context.Context, in UpdatePriceRuleInput) (*domain.PriceRule, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
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

	updated, err := u.repo.Update(ctx, in.TenantID, in.ID, r)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, apperror.New("not_found", "regla no encontrada", map[string]any{"id": in.ID.String()}, 404)
	}
	return updated, nil
}
