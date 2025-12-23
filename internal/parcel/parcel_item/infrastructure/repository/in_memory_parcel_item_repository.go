package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_item/domain"
	"ms-parcel-core/internal/parcel/parcel_item/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type InMemoryParcelItemRepository struct {
	mu   sync.Mutex
	data map[string]map[uuid.UUID]map[uuid.UUID]domain.ParcelItem // tenant -> parcel -> item -> item
}

var _ port.ParcelItemRepository = (*InMemoryParcelItemRepository)(nil)

func NewInMemoryParcelItemRepository() *InMemoryParcelItemRepository {
	return &InMemoryParcelItemRepository{data: map[string]map[uuid.UUID]map[uuid.UUID]domain.ParcelItem{}}
}

func (r *InMemoryParcelItemRepository) Add(ctx context.Context, tenantID string, item domain.ParcelItem) (uuid.UUID, error) {
	_ = ctx

	parcelID, err := uuid.Parse(item.ParcelID)
	if err != nil {
		return uuid.Nil, apperror.NewBadRequest("validation_error", "parcel_id inválido", map[string]any{"field": "parcel_id"})
	}
	itemID, err := uuid.Parse(item.ID)
	if err != nil {
		itemID = uuid.New()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return uuid.Nil, apperror.NewInternal("internal_error", "repositorio items no inicializado", nil)
	}
	if _, ok := r.data[tenantID]; !ok {
		r.data[tenantID] = map[uuid.UUID]map[uuid.UUID]domain.ParcelItem{}
	}
	if _, ok := r.data[tenantID][parcelID]; !ok {
		r.data[tenantID][parcelID] = map[uuid.UUID]domain.ParcelItem{}
	}

	item.ID = itemID.String()
	item.ParcelID = parcelID.String()
	r.data[tenantID][parcelID][itemID] = item
	return itemID, nil
}

func (r *InMemoryParcelItemRepository) ListByParcelID(ctx context.Context, tenantID string, parcelID uuid.UUID) ([]domain.ParcelItem, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio items no inicializado", nil)
	}

	byTenant, ok := r.data[tenantID]
	if !ok {
		return []domain.ParcelItem{}, nil
	}
	byParcel, ok := byTenant[parcelID]
	if !ok {
		return []domain.ParcelItem{}, nil
	}

	out := make([]domain.ParcelItem, 0, len(byParcel))
	for _, it := range byParcel {
		out = append(out, it)
	}
	return out, nil
}

func (r *InMemoryParcelItemRepository) Delete(ctx context.Context, tenantID string, parcelID uuid.UUID, itemID uuid.UUID) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return apperror.NewInternal("internal_error", "repositorio items no inicializado", nil)
	}

	byTenant, ok := r.data[tenantID]
	if !ok {
		return nil
	}
	byParcel, ok := byTenant[parcelID]
	if !ok {
		return nil
	}
	delete(byParcel, itemID)
	return nil
}
