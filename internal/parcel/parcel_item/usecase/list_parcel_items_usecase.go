package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_item/domain"
	"ms-parcel-core/internal/parcel/parcel_item/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type ListParcelItemsUseCase struct {
	parcelReader coreport.ParcelReader
	repo         port.ParcelItemRepository
}

func NewListParcelItemsUseCase(parcelReader coreport.ParcelReader, repo port.ParcelItemRepository) *ListParcelItemsUseCase {
	return &ListParcelItemsUseCase{parcelReader: parcelReader, repo: repo}
}

func (u *ListParcelItemsUseCase) Execute(ctx context.Context, tenantID string, parcelID uuid.UUID) ([]domain.ParcelItem, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if parcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	p, err := u.parcelReader.GetByID(ctx, tenantID, parcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": parcelID.String()}, 404)
	}

	return u.repo.ListByParcelID(ctx, tenantID, parcelID)
}
