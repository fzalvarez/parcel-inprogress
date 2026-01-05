package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	itemdomain "ms-parcel-core/internal/parcel/parcel_item/domain"
	itemport "ms-parcel-core/internal/parcel/parcel_item/port"
	paymentdomain "ms-parcel-core/internal/parcel/parcel_payment/domain"
	paymentport "ms-parcel-core/internal/parcel/parcel_payment/port"
	trackingdomain "ms-parcel-core/internal/parcel/parcel_tracking/domain"
	trackingport "ms-parcel-core/internal/parcel/parcel_tracking/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

const DefaultTrackingLimit = 20

type GetParcelSummaryResult struct {
	Parcel   any
	Items    []itemdomain.ParcelItem
	Payment  *paymentdomain.ParcelPayment
	Tracking []trackingdomain.TrackingEvent
}

type GetParcelSummaryUseCase struct {
	parcelRepo   coreport.ParcelReader
	itemRepo     itemport.ParcelItemRepository
	paymentRepo  paymentport.ParcelPaymentRepository
	trackingRepo trackingport.TrackingRepository
}

func NewGetParcelSummaryUseCase(parcelRepo coreport.ParcelReader, itemRepo itemport.ParcelItemRepository, paymentRepo paymentport.ParcelPaymentRepository, trackingRepo trackingport.TrackingRepository) *GetParcelSummaryUseCase {
	return &GetParcelSummaryUseCase{parcelRepo: parcelRepo, itemRepo: itemRepo, paymentRepo: paymentRepo, trackingRepo: trackingRepo}
}

func (u *GetParcelSummaryUseCase) Execute(ctx context.Context, tenantID string, parcelID uuid.UUID) (*GetParcelSummaryResult, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if parcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	p, err := u.parcelRepo.GetByID(ctx, tenantID, parcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": parcelID.String()}, 404)
	}

	items, err := u.itemRepo.ListByParcelID(ctx, tenantID, parcelID)
	if err != nil {
		return nil, err
	}

	payment, err := u.paymentRepo.GetByParcelID(ctx, tenantID, parcelID)
	if err != nil {
		return nil, err
	}

	events, err := u.trackingRepo.ListByParcelID(ctx, tenantID, parcelID.String())
	if err != nil {
		return nil, err
	}
	if len(events) > DefaultTrackingLimit {
		events = events[:DefaultTrackingLimit]
	}

	return &GetParcelSummaryResult{Parcel: p, Items: items, Payment: payment, Tracking: events}, nil
}
