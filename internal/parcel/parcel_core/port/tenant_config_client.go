package port

import "context"

type TenantConfigClient interface {
	IsEnabled(ctx context.Context, tenantID string, flagKey string) (bool, error)
}
