package usecase

import (
	"context"
	"strings"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type ListParcelsInput struct {
	TenantID string
	Filters  port.ListParcelFilters
}

type ListParcelsOutput struct {
	Items []domain.Parcel
	Count int
}

type ListParcelsUseCase struct {
	repo port.ParcelRepository
}

func NewListParcelsUseCase(repo port.ParcelRepository) *ListParcelsUseCase {
	return &ListParcelsUseCase{repo: repo}
}

func (u *ListParcelsUseCase) Execute(ctx context.Context, in ListParcelsInput) (*ListParcelsOutput, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	items, count, err := u.repo.List(ctx, in.TenantID, in.Filters)
	if err != nil {
		return nil, err
	}
	return &ListParcelsOutput{Items: items, Count: count}, nil
}
