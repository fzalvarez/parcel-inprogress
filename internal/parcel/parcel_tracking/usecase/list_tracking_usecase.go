package usecase

import (
	"context"
	"strings"

	"ms-parcel-core/internal/parcel/parcel_tracking/domain"
	"ms-parcel-core/internal/parcel/parcel_tracking/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type ListTrackingUseCase struct {
	repo port.TrackingRepository
}

func NewListTrackingUseCase(repo port.TrackingRepository) *ListTrackingUseCase {
	return &ListTrackingUseCase{repo: repo}
}

func (u *ListTrackingUseCase) Execute(ctx context.Context, tenantID string, parcelID string) ([]domain.TrackingEvent, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if strings.TrimSpace(parcelID) == "" {
		return nil, apperror.NewBadRequest("validation_error", "parcel_id inválido", map[string]any{"field": "parcel_id"})
	}
	return u.repo.ListByParcelID(ctx, tenantID, parcelID)
}
