package usecase

import (
	"context"
	"strings"

	"ms-parcel-core/internal/parcel/parcel_pricing/domain"
	"ms-parcel-core/internal/parcel/parcel_pricing/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type ListPriceRulesUseCase struct {
	repo port.PriceRuleRepository
}

func NewListPriceRulesUseCase(repo port.PriceRuleRepository) *ListPriceRulesUseCase {
	return &ListPriceRulesUseCase{repo: repo}
}

func (u *ListPriceRulesUseCase) Execute(ctx context.Context, tenantID string) ([]domain.PriceRule, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	return u.repo.List(ctx, tenantID)
}
