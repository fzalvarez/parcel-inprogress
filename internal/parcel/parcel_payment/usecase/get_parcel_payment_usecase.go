package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_payment/domain"
	"ms-parcel-core/internal/parcel/parcel_payment/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type GetParcelPaymentUseCase struct {
	repo port.ParcelPaymentRepository
}

func NewGetParcelPaymentUseCase(repo port.ParcelPaymentRepository) *GetParcelPaymentUseCase {
	return &GetParcelPaymentUseCase{repo: repo}
}

func (u *GetParcelPaymentUseCase) Execute(ctx context.Context, tenantID string, parcelID uuid.UUID) (*domain.ParcelPayment, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if parcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}
	return u.repo.GetByParcelID(ctx, tenantID, parcelID)
}
