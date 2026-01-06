package clients

import (
	"context"

	"ms-parcel-core/internal/parcel/parcel_core/port"
)

type CashboxStubClient struct{}

func NewCashboxStubClient() *CashboxStubClient {
	return &CashboxStubClient{}
}

var _ port.CashboxClient = (*CashboxStubClient)(nil)

func (c *CashboxStubClient) IsOpen(ctx context.Context, tenantID string, cashboxID string) (bool, error) {
	_ = ctx
	_ = tenantID
	_ = cashboxID
	// TODO: llamar ms-cashbox por HTTP
	return true, nil
}
