package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	docdomain "ms-parcel-core/internal/parcel/parcel_documents/domain"
	docport "ms-parcel-core/internal/parcel/parcel_documents/port"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type RegisterPrintInput struct {
	TenantID string
	ParcelID uuid.UUID
	DocType  docdomain.DocumentType
	UserID   *string
}

type RegisterPrintMeta struct {
	CountAfter        int  `json:"count_after"`
	MaxPrints         int  `json:"max_prints"`
	IsReprint         bool `json:"is_reprint"`
	ReprintFeeEnabled bool `json:"reprint_fee_enabled"`
}

type RegisterPrintResult struct {
	Record *docdomain.PrintRecord
	Meta   RegisterPrintMeta
}

type RegisterPrintUseCase struct {
	parcelRepo coreport.ParcelReader
	printRepo  docport.PrintRepository
	opts       coreport.TenantOptionsProvider
	qrGen      docport.QRGenerator
}

func NewRegisterPrintUseCase(parcelRepo coreport.ParcelReader, printRepo docport.PrintRepository, opts coreport.TenantOptionsProvider, qrGen docport.QRGenerator) *RegisterPrintUseCase {
	return &RegisterPrintUseCase{parcelRepo: parcelRepo, printRepo: printRepo, opts: opts, qrGen: qrGen}
}

func (u *RegisterPrintUseCase) Execute(ctx context.Context, in RegisterPrintInput) (*RegisterPrintResult, error) {
	if strings.TrimSpace(in.TenantID) == "" {
		return nil, apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil)
	}
	if in.ParcelID == uuid.Nil {
		return nil, apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"})
	}

	switch in.DocType {
	case docdomain.DocumentTypeLabel, docdomain.DocumentTypeReceipt, docdomain.DocumentTypeManifest, docdomain.DocumentTypeGuide:
	default:
		return nil, apperror.NewBadRequest("validation_error", "document_type inválido", map[string]any{"field": "document_type"})
	}

	p, err := u.parcelRepo.GetByID(ctx, in.TenantID, in.ParcelID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperror.New("not_found", "parcel no encontrado", map[string]any{"id": in.ParcelID.String()}, 404)
	}

	// Si es LABEL, intentar generar QR (no bloqueante)
	if in.DocType == docdomain.DocumentTypeLabel && u.qrGen != nil {
		trackingCode := ""
		if p != nil {
			trackingCode = strings.TrimSpace(p.TrackingCode)
		}
		payload := docport.QRPayload{
			TenantID:     in.TenantID,
			ParcelID:     in.ParcelID.String(),
			TrackingCode: trackingCode,
		}
		if _, err := u.qrGen.Generate(ctx, payload); err != nil {
			// TODO: logger
		}
	}

	defaults := coreport.ParcelOptions{
		RequirePackageKey:       true,
		UsePriceTable:           true,
		AllowManualPrice:        false,
		AllowOverridePriceTable: true,
		AllowPayInDestination:   false,
		MaxPrints:               1,
		AllowReprint:            false,
		ReprintFeeEnabled:       false,
	}
	opts := defaults
	if u.opts != nil {
		if o, err := u.opts.GetParcelOptions(ctx, in.TenantID); err == nil {
			opts = o
		} else {
			// TODO: logger
		}
	}
	if opts.MaxPrints <= 0 {
		opts.MaxPrints = 1
	}

	current, err := u.printRepo.CountByParcelAndType(ctx, in.TenantID, in.ParcelID, in.DocType)
	if err != nil {
		return nil, err
	}

	isReprint := current > 0
	if current >= opts.MaxPrints {
		if !opts.AllowReprint {
			return nil, apperror.New("reprint_not_allowed", "reimpresión no permitida", map[string]any{"max_prints": opts.MaxPrints, "current": current}, 409)
		}
		isReprint = true
	}

	now := time.Now().UTC()
	rec := docdomain.PrintRecord{
		ID:              uuid.NewString(),
		TenantID:        in.TenantID,
		ParcelID:        in.ParcelID.String(),
		DocumentType:    in.DocType,
		PrintedAt:       now,
		PrintedByUserID: in.UserID,
	}

	saved, err := u.printRepo.Add(ctx, in.TenantID, rec)
	if err != nil {
		return nil, err
	}

	countAfter := current + 1

	return &RegisterPrintResult{
		Record: saved,
		Meta: RegisterPrintMeta{
			CountAfter:        countAfter,
			MaxPrints:         opts.MaxPrints,
			IsReprint:         isReprint,
			ReprintFeeEnabled: opts.ReprintFeeEnabled,
		},
	}, nil
}
