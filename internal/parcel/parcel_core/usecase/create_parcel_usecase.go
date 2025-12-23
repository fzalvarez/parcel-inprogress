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
	repo         port.ParcelRepository
	tenantConfig port.TenantConfigClient
}

func NewCreateParcelUseCase(repo port.ParcelRepository, tenantConfig port.TenantConfigClient) *CreateParcelUseCase {
	return &CreateParcelUseCase{repo: repo, tenantConfig: tenantConfig}
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

	h := sha256.Sum256([]byte(in.PackageKey))
	hashHex := hex.EncodeToString(h[:])

	p := domain.Parcel{
		ID:                   uuid.NewString(),
		TenantID:             in.TenantID,
		OriginOfficeID:       in.OriginOfficeID,
		DestinationOfficeID:  in.DestinationOfficeID,
		SenderPersonID:       in.SenderPersonID,
		RecipientPersonID:    in.RecipientPersonID,
		ShipmentType:         in.ShipmentType,
		Notes:                in.Notes,
		PackageKeyHashSHA256: hashHex,
		Status:               domain.ParcelStatusCreated,
		CreatedByUserID:      in.UserID,
		CreatedByUserName:    in.UserName,
		CreatedAt:            time.Now().UTC(),
	}

	id, err := u.repo.Create(ctx, p)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
