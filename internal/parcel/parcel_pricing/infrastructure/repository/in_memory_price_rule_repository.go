package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_pricing/domain"
	"ms-parcel-core/internal/parcel/parcel_pricing/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type InMemoryPriceRuleRepository struct {
	mu   sync.Mutex
	data map[string]map[uuid.UUID]domain.PriceRule
}

var _ port.PriceRuleRepository = (*InMemoryPriceRuleRepository)(nil)

func NewInMemoryPriceRuleRepository() *InMemoryPriceRuleRepository {
	return &InMemoryPriceRuleRepository{data: map[string]map[uuid.UUID]domain.PriceRule{}}
}

func (r *InMemoryPriceRuleRepository) Create(ctx context.Context, tenantID string, rule domain.PriceRule) (*domain.PriceRule, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio pricing no inicializado", nil)
	}
	if _, ok := r.data[tenantID]; !ok {
		r.data[tenantID] = map[uuid.UUID]domain.PriceRule{}
	}

	now := time.Now().UTC()
	id := uuid.New()
	rule.ID = id.String()
	rule.TenantID = tenantID
	rule.CreatedAt = now
	rule.UpdatedAt = now

	r.data[tenantID][id] = rule
	cp := rule
	return &cp, nil
}

func (r *InMemoryPriceRuleRepository) Update(ctx context.Context, tenantID string, id uuid.UUID, rule domain.PriceRule) (*domain.PriceRule, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio pricing no inicializado", nil)
	}
	byTenant, ok := r.data[tenantID]
	if !ok {
		return nil, nil
	}
	existing, ok := byTenant[id]
	if !ok {
		return nil, nil
	}

	now := time.Now().UTC()
	rule.ID = id.String()
	rule.TenantID = tenantID
	rule.CreatedAt = existing.CreatedAt
	rule.UpdatedAt = now

	byTenant[id] = rule
	r.data[tenantID] = byTenant

	cp := rule
	return &cp, nil
}

func (r *InMemoryPriceRuleRepository) List(ctx context.Context, tenantID string) ([]domain.PriceRule, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio pricing no inicializado", nil)
	}
	byTenant, ok := r.data[tenantID]
	if !ok {
		return []domain.PriceRule{}, nil
	}

	out := make([]domain.PriceRule, 0, len(byTenant))
	for _, r := range byTenant {
		out = append(out, r)
	}
	return out, nil
}

func (r *InMemoryPriceRuleRepository) FindMatch(ctx context.Context, tenantID string, shipmentType, originOfficeID, destinationOfficeID string) (*domain.PriceRule, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio pricing no inicializado", nil)
	}
	byTenant, ok := r.data[tenantID]
	if !ok {
		return nil, nil
	}

	for _, rule := range byTenant {
		if !rule.Active {
			continue
		}
		if string(rule.ShipmentType) != shipmentType {
			continue
		}
		if rule.OriginOfficeID != originOfficeID {
			continue
		}
		if rule.DestinationOfficeID != destinationOfficeID {
			continue
		}
		cp := rule
		return &cp, nil
	}

	return nil, nil
}
