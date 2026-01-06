package port

import "context"

type QRPayload struct {
	TenantID     string
	ParcelID     string
	TrackingCode string
}

type QRGenerator interface {
	Generate(ctx context.Context, payload QRPayload) ([]byte, error)
}
