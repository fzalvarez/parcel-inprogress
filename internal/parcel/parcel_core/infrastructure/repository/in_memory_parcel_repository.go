package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type InMemoryParcelRepository struct {
	mu   sync.Mutex
	data map[uuid.UUID]domain.Parcel
}

var _ port.ParcelRepository = (*InMemoryParcelRepository)(nil)

func NewInMemoryParcelRepository() *InMemoryParcelRepository {
	return &InMemoryParcelRepository{data: map[uuid.UUID]domain.Parcel{}}
}

func (r *InMemoryParcelRepository) Create(ctx context.Context, p domain.Parcel) (uuid.UUID, error) {
	_ = ctx

	id, err := uuid.Parse(p.ID)
	if err != nil {
		id = uuid.New()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return uuid.Nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}
	r.data[id] = p
	return id, nil
}
