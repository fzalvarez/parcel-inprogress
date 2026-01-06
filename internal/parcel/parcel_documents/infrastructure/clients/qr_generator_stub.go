package clients

import (
	"context"

	"ms-parcel-core/internal/parcel/parcel_documents/port"
)

type StubQRGenerator struct{}

func NewStubQRGenerator() *StubQRGenerator {
	return &StubQRGenerator{}
}

var _ port.QRGenerator = (*StubQRGenerator)(nil)

func (s *StubQRGenerator) Generate(ctx context.Context, payload port.QRPayload) ([]byte, error) {
	_ = ctx
	_ = payload
	// TODO: llamar ms-qr-generator por HTTP y retornar bytes
	return []byte("stub-qr"), nil
}
