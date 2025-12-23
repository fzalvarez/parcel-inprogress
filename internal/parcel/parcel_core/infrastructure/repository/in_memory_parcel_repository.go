package repository

import (
	"context"
	"sort"
	"sync"
	"time"

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

func (r *InMemoryParcelRepository) GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*domain.Parcel, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	p, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	if p.TenantID != tenantID {
		// No filtrar existencia entre tenants
		return nil, nil
	}

	cp := p
	return &cp, nil
}

func (r *InMemoryParcelRepository) UpdateRegistered(ctx context.Context, tenantID string, id uuid.UUID, registeredAtUTC time.Time, userID string, userName string) (*domain.Parcel, error) {
	_ = ctx
	_ = userID
	_ = userName

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	p, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	if p.TenantID != tenantID {
		return nil, nil
	}

	p.Status = domain.ParcelStatusRegistered
	p.RegisteredAt = &registeredAtUTC
	// TODO: guardar también quién registró si se requiere en el futuro.

	r.data[id] = p
	cp := p
	return &cp, nil
}

func (r *InMemoryParcelRepository) UpdateBoarded(ctx context.Context, tenantID string, id uuid.UUID, boardedAtUTC time.Time, vehicleID string, tripID *string, departureAt *time.Time, boardedByUserID *string) (*domain.Parcel, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	p, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	if p.TenantID != tenantID {
		return nil, nil
	}

	p.Status = domain.ParcelStatusBoarded
	p.BoardedAt = &boardedAtUTC
	p.BoardedVehicleID = &vehicleID
	p.BoardedTripID = tripID
	p.BoardedDepartureAt = departureAt
	p.BoardedByUserID = boardedByUserID

	r.data[id] = p
	cp := p
	return &cp, nil
}

func (r *InMemoryParcelRepository) UpdateDelivered(ctx context.Context, tenantID string, id uuid.UUID, deliveredAtUTC time.Time, deliveredByUserID *string) (*domain.Parcel, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	p, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	if p.TenantID != tenantID {
		return nil, nil
	}

	p.Status = domain.ParcelStatusDelivered
	p.DeliveredAt = &deliveredAtUTC
	p.DeliveredByUserID = deliveredByUserID

	r.data[id] = p
	cp := p
	return &cp, nil
}

func (r *InMemoryParcelRepository) UpdateArrivedDestination(ctx context.Context, tenantID string, id uuid.UUID, arrivedAtUTC time.Time, arrivedByUserID *string) (*domain.Parcel, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	p, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	if p.TenantID != tenantID {
		return nil, nil
	}

	p.Status = domain.ParcelStatusArrivedDestination
	p.ArrivedAt = &arrivedAtUTC
	p.ArrivedByUserID = arrivedByUserID

	r.data[id] = p
	cp := p
	return &cp, nil
}

func (r *InMemoryParcelRepository) UpdateInTransit(ctx context.Context, tenantID string, id uuid.UUID, departedAtUTC time.Time, departedByUserID *string, vehicleID *string) (*domain.Parcel, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	p, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	if p.TenantID != tenantID {
		return nil, nil
	}

	p.Status = domain.ParcelStatusInTransit
	p.DepartedAt = &departedAtUTC
	p.DepartedByUserID = departedByUserID

	if vehicleID != nil {
		// En MVP usamos boarded_vehicle_id como referencia del vehículo de tránsito
		p.BoardedVehicleID = vehicleID
	}

	r.data[id] = p
	cp := p
	return &cp, nil
}

func (r *InMemoryParcelRepository) ListByFilters(ctx context.Context, tenantID string, f port.ListParcelFilters) ([]domain.Parcel, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	out := make([]domain.Parcel, 0)
	for _, p := range r.data {
		if p.TenantID != tenantID {
			continue
		}
		if f.Status != nil && p.Status != *f.Status {
			continue
		}
		if f.VehicleID != nil {
			if p.BoardedVehicleID == nil || *p.BoardedVehicleID != *f.VehicleID {
				continue
			}
		}
		if f.OriginOfficeID != nil && p.OriginOfficeID != *f.OriginOfficeID {
			continue
		}
		if f.DestinationOfficeID != nil && p.DestinationOfficeID != *f.DestinationOfficeID {
			continue
		}

		out = append(out, p)
	}

	return out, nil
}

func (r *InMemoryParcelRepository) List(ctx context.Context, tenantID string, f port.ListParcelFilters) ([]domain.Parcel, int, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, 0, apperror.NewInternal("internal_error", "repositorio no inicializado", nil)
	}

	filtered := make([]domain.Parcel, 0)
	for _, p := range r.data {
		if p.TenantID != tenantID {
			continue
		}
		if f.Status != nil && p.Status != *f.Status {
			continue
		}
		if f.VehicleID != nil {
			if p.BoardedVehicleID == nil || *p.BoardedVehicleID != *f.VehicleID {
				continue
			}
		}
		if f.OriginOfficeID != nil && p.OriginOfficeID != *f.OriginOfficeID {
			continue
		}
		if f.DestinationOfficeID != nil && p.DestinationOfficeID != *f.DestinationOfficeID {
			continue
		}
		if f.SenderPersonID != nil && p.SenderPersonID != *f.SenderPersonID {
			continue
		}
		if f.RecipientPersonID != nil && p.RecipientPersonID != *f.RecipientPersonID {
			continue
		}
		if f.FromCreatedAt != nil && p.CreatedAt.Before(f.FromCreatedAt.UTC()) {
			continue
		}
		if f.ToCreatedAt != nil && p.CreatedAt.After(f.ToCreatedAt.UTC()) {
			continue
		}

		filtered = append(filtered, p)
	}

	// Orden created_at desc
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	count := len(filtered)
	limit := f.Limit
	offset := f.Offset
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	start := offset
	if start > count {
		start = count
	}
	end := start + limit
	if end > count {
		end = count
	}

	paged := make([]domain.Parcel, 0, end-start)
	paged = append(paged, filtered[start:end]...)

	return paged, count, nil
}
