package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"

	"ms-parcel-core/internal/parcel/parcel_core/domain"
	repository "ms-parcel-core/internal/parcel/parcel_core/infrastructure/repository"
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

func buildYearCode(now time.Time) string {
	y := now.UTC().Year()
	if y < 2025 {
		return "A"
	}
	offset := y - 2025
	if offset < 0 {
		offset = 0
	}
	if offset > 25 {
		offset = 25
	}
	return string(rune('A' + offset))
}

func randomCrockford(n int) (string, error) {
	const alphabet = "ABCDEFGHJKMNPQRSTVWXYZ23456789"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(out), nil
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

	if strings.TrimSpace(p.TrackingCode) == "" {
		now := time.Now().UTC()
		assigned := false

		checker, ok := u.repo.(*repository.InMemoryParcelRepository)
		if !ok {
			return uuid.Nil, apperror.NewInternal("internal_error", "repositorio no soporta verificación de tracking_code", nil)
		}

		for i := 0; i < 5; i++ {
			rnd, err := randomCrockford(5)
			if err != nil {
				return uuid.Nil, apperror.NewInternal("internal_error", "no se pudo generar tracking_code", map[string]any{"error": err.Error()})
			}
			code := "QB" + buildYearCode(now) + rnd
			exists, err := checker.ExistsTrackingCode(ctx, code)
			if err != nil {
				return uuid.Nil, err
			}
			if !exists {
				p.TrackingCode = code
				assigned = true
				break
			}
		}
		if !assigned {
			return uuid.Nil, apperror.NewInternal("internal_error", "no se pudo asignar tracking_code", nil)
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
