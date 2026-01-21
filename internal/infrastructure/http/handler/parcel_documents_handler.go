package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	docdomain "ms-parcel-core/internal/parcel/parcel_documents/domain"
	docport "ms-parcel-core/internal/parcel/parcel_documents/port"
	docusecase "ms-parcel-core/internal/parcel/parcel_documents/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type RegisterPrintRequest struct {
	DocumentType string `json:"document_type" binding:"required,oneof=LABEL RECEIPT MANIFEST GUIDE"`
}

type PrintRecordResponse struct {
	ID              string  `json:"id"`
	ParcelID        string  `json:"parcel_id"`
	DocumentType    string  `json:"document_type"`
	PrintedAt       string  `json:"printed_at"`
	PrintedByUserID *string `json:"printed_by_user_id,omitempty"`
}

type RegisterPrintResponse struct {
	Record PrintRecordResponse          `json:"record"`
	Meta   docusecase.RegisterPrintMeta `json:"meta"`
}

type ParcelDocumentsHandler struct {
	registerUC *docusecase.RegisterPrintUseCase
	printRepo  docport.PrintRepository
}

func NewParcelDocumentsHandler(registerUC *docusecase.RegisterPrintUseCase, printRepo docport.PrintRepository) *ParcelDocumentsHandler {
	return &ParcelDocumentsHandler{registerUC: registerUC, printRepo: printRepo}
}

// RegisterPrint godoc
// @Summary Registrar impresión de documento
// @Description Registra un evento de impresión de documento (LABEL, RECEIPT, MANIFEST, GUIDE) para un envío. El registro incluye tipo de documento, timestamp y usuario que realizó la impresión.
// @Tags ParcelDocuments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Param payload body RegisterPrintRequest true "Solicitud de impresión con tipo de documento"
// @Success 200 {object} handler.AnyDataEnvelope "Impresión registrada exitosamente"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido, payload malformado o tipo de documento no permitido"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 404 {object} handler.ErrorResponse "Envío no encontrado"
// @Failure 409 {object} handler.ErrorResponse "Conflicto: estado incompatible o límite de impresiones alcanzado"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /parcels/{id}/documents/print [post]
func (h *ParcelDocumentsHandler) RegisterPrint(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	parcelID, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	var req RegisterPrintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	userIDVal, _ := c.Get("user_id")
	uid := strings.TrimSpace(anyToString(userIDVal))
	var uidPtr *string
	if uid != "" {
		uidPtr = &uid
	}

	res, err := h.registerUC.Execute(c.Request.Context(), docusecase.RegisterPrintInput{
		TenantID: tenant,
		ParcelID: parcelID,
		DocType:  docdomain.DocumentType(strings.TrimSpace(req.DocumentType)),
		UserID:   uidPtr,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": RegisterPrintResponse{
			Record: PrintRecordResponse{
				ID:              res.Record.ID,
				ParcelID:        res.Record.ParcelID,
				DocumentType:    string(res.Record.DocumentType),
				PrintedAt:       res.Record.PrintedAt.UTC().Format(time.RFC3339),
				PrintedByUserID: res.Record.PrintedByUserID,
			},
			Meta: res.Meta,
		},
	})
}

// ListPrints godoc
// @Summary Listar impresiones de envío
// @Description Lista todos los registros de impresión asociados a un envío específico. Incluye información de timestamp, tipo de documento e usuario que realizó la impresión.
// @Tags ParcelDocuments
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID del envío" Format(uuid)
// @Success 200 {object} handler.AnyDataEnvelope "Lista de registros de impresión"
// @Failure 400 {object} handler.ErrorResponse "Validación fallida: id inválido"
// @Failure 401 {object} handler.ErrorResponse "No autorizado: token inválido o credenciales faltantes"
// @Failure 500 {object} handler.ErrorResponse "Error interno del servidor"
// @Router /parcels/{id}/documents/prints [get]
func (h *ParcelDocumentsHandler) ListPrints(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	parcelID, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	recs, err := h.printRepo.ListByParcel(c.Request.Context(), tenant, parcelID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	out := make([]PrintRecordResponse, 0, len(recs))
	for _, r := range recs {
		out = append(out, PrintRecordResponse{
			ID:              r.ID,
			ParcelID:        r.ParcelID,
			DocumentType:    string(r.DocumentType),
			PrintedAt:       r.PrintedAt.UTC().Format(time.RFC3339),
			PrintedByUserID: r.PrintedByUserID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
}
