package port

import "context"

type CashboxClient interface {
	IsOpen(ctx context.Context, tenantID string, cashboxID string) (bool, error)
}
