package repository

import (
	"context"
	"sort"
	"sync"

	"ms-parcel-core/internal/parcel/parcel_tracking/domain"
	"ms-parcel-core/internal/parcel/parcel_tracking/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type InMemoryTrackingRepository struct {
	mu   sync.Mutex
	data map[string]map[string][]domain.TrackingEvent // tenantID -> parcelID -> events
}

var _ port.TrackingRepository = (*InMemoryTrackingRepository)(nil)

func NewInMemoryTrackingRepository() *InMemoryTrackingRepository {
	return &InMemoryTrackingRepository{data: map[string]map[string][]domain.TrackingEvent{}}
}

func (r *InMemoryTrackingRepository) Append(ctx context.Context, tenantID string, ev domain.TrackingEvent) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return apperror.NewInternal("internal_error", "repositorio tracking no inicializado", nil)
	}
	if _, ok := r.data[tenantID]; !ok {
		r.data[tenantID] = map[string][]domain.TrackingEvent{}
	}
	r.data[tenantID][ev.ParcelID] = append(r.data[tenantID][ev.ParcelID], ev)
	return nil
}

func (r *InMemoryTrackingRepository) ListByParcelID(ctx context.Context, tenantID string, parcelID string) ([]domain.TrackingEvent, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio tracking no inicializado", nil)
	}

	byTenant, ok := r.data[tenantID]
	if !ok {
		return []domain.TrackingEvent{}, nil
	}
	evs := byTenant[parcelID]

	out := make([]domain.TrackingEvent, 0, len(evs))
	out = append(out, evs...)

	sort.Slice(out, func(i, j int) bool {
		return out[i].OccurredAt.Before(out[j].OccurredAt)
	})

	return out, nil
}
