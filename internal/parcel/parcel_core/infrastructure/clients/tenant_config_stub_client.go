package clients

import "context"

type TenantConfigStubClient struct{}

func NewTenantConfigStubClient() *TenantConfigStubClient {
	return &TenantConfigStubClient{}
}

func (c *TenantConfigStubClient) IsEnabled(ctx context.Context, tenantID string, flagKey string) (bool, error) {
	_ = ctx
	_ = tenantID
	_ = flagKey
	// TODO: integrar cliente real a TENANT-CONFIG.
	return false, nil
}
