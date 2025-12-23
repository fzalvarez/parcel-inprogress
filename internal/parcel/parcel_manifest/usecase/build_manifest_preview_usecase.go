package usecase

import (
	"context"
	"strings"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	manifestdomain "ms-parcel-core/internal/parcel/parcel_manifest/domain"
	manifestport "ms-parcel-core/internal/parcel/parcel_manifest/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type BuildManifestPreviewInput struct {
	TenantID            string
	VehicleID           string
	OriginOfficeID      string
	DestinationOfficeID string
}

type BuildManifestPreviewUseCase struct {
	reader manifestport.ParcelReader
}

func NewBuildManifestPreviewUseCase(reader manifestport.ParcelReader) *BuildManifestPreviewUseCase {
	return &BuildManifestPreviewUseCase{reader: reader}
}

func (u *BuildManifestPreviewUseCase) Execute(ctx context.Context, in BuildManifestPreviewInput) (*manifestdomain.ManifestPreview, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}

	status := domain.ParcelStatusBoarded
	vid := in.VehicleID
	oid := in.OriginOfficeID
	did := in.DestinationOfficeID

	parcels, err := u.reader.ListByFilters(ctx, in.TenantID, coreport.ListParcelFilters{
		Status:              &status,
		VehicleID:           &vid,
		OriginOfficeID:      &oid,
		DestinationOfficeID: &did,
	})
	if err != nil {
		return nil, err
	}

	prev := &manifestdomain.ManifestPreview{
		VehicleID:           in.VehicleID,
		OriginOfficeID:      in.OriginOfficeID,
		DestinationOfficeID: in.DestinationOfficeID,
		Parcels:             make([]manifestdomain.ParcelSummary, 0, len(parcels)),
	}

	for _, p := range parcels {
		prev.Parcels = append(prev.Parcels, manifestdomain.ParcelSummary{
			ParcelID:          p.ID,
			Status:            string(p.Status),
			SenderPersonID:    p.SenderPersonID,
			RecipientPersonID: p.RecipientPersonID,
			Notes:             p.Notes,
		})
	}

	prev.Totals = manifestdomain.ManifestTotals{CountParcels: len(prev.Parcels)}
	return prev, nil
}
