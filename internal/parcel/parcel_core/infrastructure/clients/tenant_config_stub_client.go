package clients

import (
	"context"

	"ms-parcel-core/internal/parcel/parcel_core/port"
)

type TenantConfigStubClient struct{}

func NewTenantConfigStubClient() *TenantConfigStubClient {
	return &TenantConfigStubClient{}
}

var _ port.TenantOptionsProvider = (*TenantConfigStubClient)(nil)
var _ port.TenantConfigClient = (*TenantConfigStubClient)(nil)

func (c *TenantConfigStubClient) IsEnabled(ctx context.Context, tenantID string, featureKey string) (bool, error) {
	_ = ctx
	_ = tenantID
	_ = featureKey
	return false, nil
}

func (c *TenantConfigStubClient) GetParcelOptions(ctx context.Context, tenantID string) (port.ParcelOptions, error) {
	_ = ctx
	_ = tenantID

	return port.ParcelOptions{
		RequirePackageKey:       true,
		UsePriceTable:           true,
		AllowManualPrice:        false,
		AllowOverridePriceTable: true,
		AllowPayInDestination:   false,
		MaxPrints:               1,
		AllowReprint:            false,
		ReprintFeeEnabled:       false,
	}, nil
}
