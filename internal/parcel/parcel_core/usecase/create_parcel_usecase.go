package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type CreateParcelInput struct {
	TenantID            string
	UserID              string
	UserName            string
	ShipmentType        domain.ShipmentType
	OriginOfficeID      string
	DestinationOfficeID string
	SenderPersonID      string
	RecipientPersonID   string
	Notes               *string
	PackageKey          string
	PackageKeyConfirm   string
}

type CreateParcelUseCase struct {
	repo            port.ParcelRepository
	tenantConfig    port.TenantConfigClient
	tracking        port.TrackingRecorder
	optionsProvider port.TenantOptionsProvider
}

func NewCreateParcelUseCase(repo port.ParcelRepository, tenantConfig port.TenantConfigClient, tracking port.TrackingRecorder, optionsProvider port.TenantOptionsProvider) *CreateParcelUseCase {
	return &CreateParcelUseCase{repo: repo, tenantConfig: tenantConfig, tracking: tracking, optionsProvider: optionsProvider}
}

func (u *CreateParcelUseCase) Execute(ctx context.Context, in CreateParcelInput) (uuid.UUID, error) {
	if strings.TrimSpace(in.PackageKey) == "" || strings.TrimSpace(in.PackageKeyConfirm) == "" {
		return uuid.Nil, apperror.NewBadRequest("validation_error", "package_key y package_key_confirm son requeridos", map[string]any{"field": "package_key"})
	}
	if in.PackageKey != in.PackageKeyConfirm {
		return uuid.Nil, apperror.NewBadRequest("validation_error", "package_key y package_key_confirm no coinciden", map[string]any{"field": "package_key_confirm"})
	}
	if strings.TrimSpace(in.TenantID) == "" || strings.TrimSpace(in.UserID) == "" {
		return uuid.Nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}

	// TODO: consultar flags reales desde TENANT-CONFIG si aplica para este flujo.
	_, _ = u.tenantConfig.IsEnabled(ctx, in.TenantID, "parcel_core.create")

	defaults := port.ParcelOptions{
		RequirePackageKey:       true,
		UsePriceTable:           true,
		AllowManualPrice:        false,
		AllowOverridePriceTable: true,
		AllowPayInDestination:   false,
	}
	opts := defaults
	if u.optionsProvider != nil {
		if o, err := u.optionsProvider.GetParcelOptions(ctx, in.TenantID); err == nil {
			opts = o
		} else {
			// TODO: logger
		}
	}

	p := domain.Parcel{
		ID:                  uuid.NewString(),
		TenantID:            in.TenantID,
		OriginOfficeID:      in.OriginOfficeID,
		DestinationOfficeID: in.DestinationOfficeID,
		SenderPersonID:      in.SenderPersonID,
		RecipientPersonID:   in.RecipientPersonID,
		ShipmentType:        in.ShipmentType,
		Notes:               in.Notes,
		Status:              domain.ParcelStatusCreated,
		CreatedByUserID:     in.UserID,
		CreatedByUserName:   in.UserName,
		CreatedAt:           time.Now().UTC(),
	}

	// Validación package_key/confirmación existente, condicionada por flag
	if opts.RequirePackageKey {
		if strings.TrimSpace(in.PackageKey) == "" {
			return uuid.Nil, apperror.NewBadRequest("validation_error", "package_key requerido", map[string]any{"field": "package_key"})
		}
		h := sha256.Sum256([]byte(in.PackageKey))
		p.PackageKeyHashSHA256 = hex.EncodeToString(h[:])
	} else {
		// Si no requiere package_key, permitir vacío y solo hashear si viene
		if strings.TrimSpace(in.PackageKey) != "" {
			h := sha256.Sum256([]byte(in.PackageKey))
			p.PackageKeyHashSHA256 = hex.EncodeToString(h[:])
		} else {
			p.PackageKeyHashSHA256 = ""
		}
	}

	id, err := u.repo.Create(ctx, p)
	if err != nil {
		return uuid.Nil, err
	}

	if u.tracking != nil {
		if err := u.tracking.RecordEvent(ctx, in.TenantID, port.TrackingEventDTO{
			ParcelID:   id.String(),
			EventType:  port.EventTypeParcelCreated,
			OccurredAt: time.Now().UTC(),
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata: map[string]any{
				"shipment_type":         string(in.ShipmentType),
				"origin_office_id":      in.OriginOfficeID,
				"destination_office_id": in.DestinationOfficeID,
			},
		}); err != nil {
			// TODO: logger
		}
	}

	return id, nil
}
