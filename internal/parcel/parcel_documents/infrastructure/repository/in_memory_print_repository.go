package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_documents/domain"
	"ms-parcel-core/internal/parcel/parcel_documents/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type InMemoryPrintRepository struct {
	mu   sync.Mutex
	data map[string]map[uuid.UUID][]domain.PrintRecord
}

var _ port.PrintRepository = (*InMemoryPrintRepository)(nil)

func NewInMemoryPrintRepository() *InMemoryPrintRepository {
	return &InMemoryPrintRepository{data: map[string]map[uuid.UUID][]domain.PrintRecord{}}
}

func (r *InMemoryPrintRepository) Add(ctx context.Context, tenantID string, rec domain.PrintRecord) (*domain.PrintRecord, error) {
	_ = ctx

	parcelID, err := uuid.Parse(rec.ParcelID)
	if err != nil {
		return nil, apperror.NewBadRequest("validation_error", "parcel_id inválido", map[string]any{"field": "parcel_id"})
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio documentos no inicializado", nil)
	}
	if _, ok := r.data[tenantID]; !ok {
		r.data[tenantID] = map[uuid.UUID][]domain.PrintRecord{}
	}

	rec.ParcelID = parcelID.String()
	r.data[tenantID][parcelID] = append(r.data[tenantID][parcelID], rec)

	cp := rec
	return &cp, nil
}

func (r *InMemoryPrintRepository) CountByParcelAndType(ctx context.Context, tenantID string, parcelID uuid.UUID, docType domain.DocumentType) (int, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return 0, apperror.NewInternal("internal_error", "repositorio documentos no inicializado", nil)
	}

	byTenant, ok := r.data[tenantID]
	if !ok {
		return 0, nil
	}
	recs, ok := byTenant[parcelID]
	if !ok {
		return 0, nil
	}

	cnt := 0
	for _, r := range recs {
		if r.DocumentType == docType {
			cnt++
		}
	}
	return cnt, nil
}

func (r *InMemoryPrintRepository) ListByParcel(ctx context.Context, tenantID string, parcelID uuid.UUID) ([]domain.PrintRecord, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		return nil, apperror.NewInternal("internal_error", "repositorio documentos no inicializado", nil)
	}

	byTenant, ok := r.data[tenantID]
	if !ok {
		return []domain.PrintRecord{}, nil
	}
	recs, ok := byTenant[parcelID]
	if !ok {
		return []domain.PrintRecord{}, nil
	}

	out := make([]domain.PrintRecord, 0, len(recs))
	out = append(out, recs...)
	return out, nil
}
