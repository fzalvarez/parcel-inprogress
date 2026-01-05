package port

import (
	"context"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_documents/domain"
)

type PrintRepository interface {
	Add(ctx context.Context, tenantID string, r domain.PrintRecord) (*domain.PrintRecord, error)
	CountByParcelAndType(ctx context.Context, tenantID string, parcelID uuid.UUID, docType domain.DocumentType) (int, error)
	ListByParcel(ctx context.Context, tenantID string, parcelID uuid.UUID) ([]domain.PrintRecord, error)
}
