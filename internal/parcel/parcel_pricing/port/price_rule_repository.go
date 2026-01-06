package port

import (
	"context"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_pricing/domain"
)

type PriceRuleRepository interface {
	Create(ctx context.Context, tenantID string, r domain.PriceRule) (*domain.PriceRule, error)
	Update(ctx context.Context, tenantID string, id uuid.UUID, r domain.PriceRule) (*domain.PriceRule, error)
	List(ctx context.Context, tenantID string) ([]domain.PriceRule, error)
	FindMatch(ctx context.Context, tenantID string, shipmentType, originOfficeID, destinationOfficeID string) (*domain.PriceRule, error)
}
