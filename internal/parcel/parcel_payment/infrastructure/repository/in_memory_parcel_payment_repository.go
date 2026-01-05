package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_payment/domain"
	"ms-parcel-core/internal/parcel/parcel_payment/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type InMemoryParcelPaymentRepository struct {
	mu   sync.Mutex
	data map[string]map[uuid.UUID]domain.ParcelPayment // tenant -> parcel -> payment
}

var _ port.ParcelPaymentRepository = (*InMemoryParcelPaymentRepository)(nil)

func NewInMemoryParcelPaymentRepository() *InMemoryParcelPaymentRepository {
	return &InMemoryParcelPaymentRepository{data: map[string]map[uuid.UUID]domain.ParcelPayment{}}
}

func (r *InMemoryParcelPaymentRepository) Upsert(ctx context.Context, tenantID string, p domain.ParcelPayment) (*domain.ParcelPayment, error) {
	_ = ctx

	parcelID, err := uuid.Parse(p.ParcelID)
	if err != nil {
		return nil, apperror.NewBadRequest("validation_error", "parcel_id inválido", map[string]any{"field": "parcel_id"})
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio pagos no inicializado", nil)
	}
	if _, ok := r.data[tenantID]; !ok {
		r.data[tenantID] = map[uuid.UUID]domain.ParcelPayment{}
	}

	r.data[tenantID][parcelID] = p
	cp := p
	return &cp, nil
}

func (r *InMemoryParcelPaymentRepository) GetByParcelID(ctx context.Context, tenantID string, parcelID uuid.UUID) (*domain.ParcelPayment, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio pagos no inicializado", nil)
	}

	byTenant, ok := r.data[tenantID]
	if !ok {
		return nil, nil
	}
	p, ok := byTenant[parcelID]
	if !ok {
		return nil, nil
	}

	cp := p
	return &cp, nil
}
