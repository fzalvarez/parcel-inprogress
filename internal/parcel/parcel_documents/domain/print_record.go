package domain

import "time"

type DocumentType string

const (
	DocumentTypeLabel    DocumentType = "LABEL"
	DocumentTypeReceipt  DocumentType = "RECEIPT"
	DocumentTypeManifest DocumentType = "MANIFEST"
	DocumentTypeGuide    DocumentType = "GUIDE"
)

type PrintRecord struct {
	ID              string
	TenantID        string
	ParcelID        string
	DocumentType    DocumentType
	PrintedAt       time.Time
	PrintedByUserID *string
}
