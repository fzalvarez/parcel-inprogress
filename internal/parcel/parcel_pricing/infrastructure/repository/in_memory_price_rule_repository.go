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

	// Búsqueda jerárquica por prioridad:
	// 1. Coincidencia exacta: Origin -> Destination
	// 2. Origen específico, destino comodín: Origin -> *
	// 3. Origen comodín, destino específico: * -> Destination
	// 4. Comodín total: * -> *

	var candidates []domain.PriceRule

	for _, rule := range byTenant {
		if !rule.Active {
			continue
		}
		if string(rule.ShipmentType) != shipmentType {
			continue
		}

		originMatch := rule.OriginOfficeID == originOfficeID || rule.OriginOfficeID == domain.WildcardOffice
		destMatch := rule.DestinationOfficeID == destinationOfficeID || rule.DestinationOfficeID == domain.WildcardOffice

		if originMatch && destMatch {
			candidates = append(candidates, rule)
		}
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	// Ordenar por especificidad (prioridad implícita)
	best := candidates[0]
	bestScore := calculateRuleScore(best, originOfficeID, destinationOfficeID)

	for i := 1; i < len(candidates); i++ {
		score := calculateRuleScore(candidates[i], originOfficeID, destinationOfficeID)
		if score > bestScore || (score == bestScore && candidates[i].Priority > best.Priority) {
			best = candidates[i]
			bestScore = score
		}
	}

	cp := best
	return &cp, nil
}

// calculateRuleScore asigna puntaje de especificidad
// Mayor puntaje = más específica = mayor prioridad
func calculateRuleScore(rule domain.PriceRule, targetOrigin, targetDest string) int {
	score := 0

	if rule.OriginOfficeID == targetOrigin {
		score += 10 // Origen exacto
	} else if rule.OriginOfficeID == domain.WildcardOffice {
		score += 1 // Origen comodín
	}

	if rule.DestinationOfficeID == targetDest {
		score += 10 // Destino exacto
	} else if rule.DestinationOfficeID == domain.WildcardOffice {
		score += 1 // Destino comodín
	}

	return score
}
