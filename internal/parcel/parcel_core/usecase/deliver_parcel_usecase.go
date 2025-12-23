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

type DeliverParcelInput struct {
	TenantID   string
	UserID     string
	UserName   string
	ParcelID   uuid.UUID
	PackageKey string
}

type DeliverParcelUseCase struct {
	repo     port.ParcelRepository
	tracking port.TrackingRecorder
}

func NewDeliverParcelUseCase(repo port.ParcelRepository, tracking port.TrackingRecorder) *DeliverParcelUseCase {
	return &DeliverParcelUseCase{repo: repo, tracking: tracking}
}

func (u *DeliverParcelUseCase) Execute(ctx context.Context, in DeliverParcelInput) (*domain.Parcel, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if strings.TrimSpace(in.UserID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}
	if strings.TrimSpace(in.PackageKey) == "" {
		return nil, apperror.NewBadRequest("validation_error", "package_key requerido", map[string]any{"field": "package_key"})
	}

	p, err := u.repo.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	// Estados entregables (por ahora)
	deliverable := p.Status == domain.ParcelStatusRegistered || p.Status == domain.ParcelStatusBoarded
	if !deliverable {
		return nil, apperror.New("invalid_state", "transición de estado inválida", map[string]any{"allowed": []domain.ParcelStatus{domain.ParcelStatusRegistered, domain.ParcelStatusBoarded}, "actual": p.Status}, 409)
	}

	h := sha256.Sum256([]byte(in.PackageKey))
	hashHex := hex.EncodeToString(h[:])
	if p.PackageKeyHashSHA256 == "" || hashHex != p.PackageKeyHashSHA256 {
		return nil, apperror.New("invalid_package_key", "package_key inválido", nil, 403)
	}

	deliveredAt := time.Now().UTC()
	by := strings.TrimSpace(in.UserID)
	byPtr := &by

	updated, err := u.repo.UpdateDelivered(ctx, in.TenantID, in.ParcelID, deliveredAt, byPtr)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	if u.tracking != nil {
		if err := u.tracking.RecordEvent(ctx, in.TenantID, port.TrackingEventDTO{
			ParcelID:   in.ParcelID.String(),
			EventType:  port.EventTypeParcelDelivered,
			OccurredAt: deliveredAt,
			UserID:     in.UserID,
			UserName:   in.UserName,
			Metadata: map[string]any{
				"delivered_at":         deliveredAt.UTC().Format(time.RFC3339),
				"delivered_by_user_id": by,
			},
		}); err != nil {
			// TODO: logger
		}
	}

	return updated, nil
}
