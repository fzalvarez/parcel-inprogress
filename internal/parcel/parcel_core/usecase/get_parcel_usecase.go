package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type GetParcelInput struct {
	TenantID string
	ParcelID uuid.UUID
}

type GetParcelUseCase struct {
	repo port.ParcelRepository
}

func NewGetParcelUseCase(repo port.ParcelRepository) *GetParcelUseCase {
	return &GetParcelUseCase{repo: repo}
}

func (u *GetParcelUseCase) Execute(ctx context.Context, in GetParcelInput) (*domain.Parcel, error) {
	if strings.TrimSpace(in.TenantID) == "" {
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
	return p, nil
}
