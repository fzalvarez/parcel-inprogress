package port

import "context"

type ParcelOptions struct {
	RequirePackageKey       bool
	UsePriceTable           bool
	AllowManualPrice        bool
	AllowOverridePriceTable bool
	AllowPayInDestination   bool
	MaxPrints               int
	AllowReprint            bool
	ReprintFeeEnabled       bool
}

type TenantOptionsProvider interface {
	GetParcelOptions(ctx context.Context, tenantID string) (ParcelOptions, error)
}
