package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ms-parcel-core/internal/infrastructure/http/dto"
	"ms-parcel-core/internal/parcel/parcel_core/domain"
	"ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_core/usecase"
	"ms-parcel-core/internal/pkg/util/apperror"
)

type ParcelHandler struct {
	createUC   *usecase.CreateParcelUseCase
	getUC      *usecase.GetParcelUseCase
	listUC     *usecase.ListParcelsUseCase
	registerUC *usecase.RegisterParcelUseCase
	boardUC    *usecase.BoardParcelUseCase
	departUC   *usecase.DepartParcelUseCase
	arriveUC   *usecase.ArriveParcelUseCase
	deliverUC  *usecase.DeliverParcelUseCase
}

func NewParcelHandler(
	createUC *usecase.CreateParcelUseCase,
	listUC *usecase.ListParcelsUseCase,
	getUC *usecase.GetParcelUseCase,
	registerUC *usecase.RegisterParcelUseCase,
	boardUC *usecase.BoardParcelUseCase,
	departUC *usecase.DepartParcelUseCase,
	arriveUC *usecase.ArriveParcelUseCase,
	deliverUC *usecase.DeliverParcelUseCase) *ParcelHandler {
	return &ParcelHandler{
		createUC:   createUC,
		listUC:     listUC,
		getUC:      getUC,
		registerUC: registerUC,
		boardUC:    boardUC,
		departUC:   departUC,
		arriveUC:   arriveUC,
		deliverUC:  deliverUC,
	}
}

// Create godoc
// @Summary Crear envío
// @Description Crea un nuevo envío (parcel)
// @Tags Parcels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param payload body CreateParcelRequest true "Create parcel"
// @Success 201 {object} handler.CreateParcelResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels [post]
func (h *ParcelHandler) Create(c *gin.Context) {
	var req dto.CreateParcelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	in := usecase.CreateParcelInput{
		TenantID:            strings.TrimSpace(anyToString(tenantID)),
		UserID:              strings.TrimSpace(anyToString(userID)),
		UserName:            strings.TrimSpace(anyToString(userName)),
		ShipmentType:        domain.ShipmentType(req.ShipmentType),
		OriginOfficeID:      req.OriginOfficeID,
		DestinationOfficeID: req.DestinationOfficeID,
		SenderPersonID:      req.SenderPersonID,
		RecipientPersonID:   req.RecipientPersonID,
		Notes:               req.Notes,
		PackageKey:          req.PackageKey,
		PackageKeyConfirm:   req.PackageKeyConfirm,
	}

	id, err := h.createUC.Execute(c.Request.Context(), in)
	if err != nil {
		_ = c.Error(err)
		return
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)

	// Header temporal de debug (remover luego de estabilizar)
	c.Header("X-Debug-Parcel", "created")

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": dto.CreateParcelResponse{
			ID:                  id.String(),
			Status:              string(domain.ParcelStatusCreated),
			ShipmentType:        req.ShipmentType,
			OriginOfficeID:      req.OriginOfficeID,
			DestinationOfficeID: req.DestinationOfficeID,
			SenderPersonID:      req.SenderPersonID,
			RecipientPersonID:   req.RecipientPersonID,
			Notes:               req.Notes,
			CreatedAt:           createdAt,
		},
	})
}

// List godoc
// @Summary Listar envíos
// @Description Lista envíos con filtros y paginación
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param q query string false "Search by tracking_code or parcel id"
// @Param status query string false "Status"
// @Param origin_office_id query string false "Origin office id"
// @Param destination_office_id query string false "Destination office id"
// @Param sender_person_id query string false "Sender person id"
// @Param recipient_person_id query string false "Recipient person id"
// @Param vehicle_id query string false "Vehicle id"
// @Param from_created_at query string false "From created_at (RFC3339)"
// @Param to_created_at query string false "To created_at (RFC3339)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} handler.ParcelListResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels [get]
func (h *ParcelHandler) List(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	var statusPtr *domain.ParcelStatus
	statusStr := strings.TrimSpace(c.Query("status"))
	if statusStr != "" {
		s := domain.ParcelStatus(statusStr)
		statusPtr = &s
	}

	parseUUIDQuery := func(key string) (*string, error) {
		v := strings.TrimSpace(c.Query(key))
		if v == "" {
			return nil, nil
		}
		if _, err := uuid.Parse(v); err != nil {
			return nil, apperror.NewBadRequest("validation_error", key+" inválido", map[string]any{"field": key})
		}
		return &v, nil
	}

	originOfficeID, err := parseUUIDQuery("origin_office_id")
	if err != nil {
		_ = c.Error(err)
		return
	}
	destinationOfficeID, err := parseUUIDQuery("destination_office_id")
	if err != nil {
		_ = c.Error(err)
		return
	}
	vehicleID, err := parseUUIDQuery("vehicle_id")
	if err != nil {
		_ = c.Error(err)
		return
	}
	senderPersonID, err := parseUUIDQuery("sender_person_id")
	if err != nil {
		_ = c.Error(err)
		return
	}
	recipientPersonID, err := parseUUIDQuery("recipient_person_id")
	if err != nil {
		_ = c.Error(err)
		return
	}

	var fromPtr *time.Time
	fromStr := strings.TrimSpace(c.Query("from"))
	if fromStr != "" {
		tm, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "from inválido", map[string]any{"field": "from"}))
			return
		}
		ut := tm.UTC()
		fromPtr = &ut
	}

	var toPtr *time.Time
	toStr := strings.TrimSpace(c.Query("to"))
	if toStr != "" {
		tm, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "to inválido", map[string]any{"field": "to"}))
			return
		}
		ut := tm.UTC()
		toPtr = &ut
	}

	limit := 50
	if l := strings.TrimSpace(c.Query("limit")); l != "" {
		v, err := strconv.Atoi(l)
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "limit inválido", map[string]any{"field": "limit"}))
			return
		}
		limit = v
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	offset := 0
	if o := strings.TrimSpace(c.Query("offset")); o != "" {
		v, err := strconv.Atoi(o)
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "offset inválido", map[string]any{"field": "offset"}))
			return
		}
		offset = v
	}
	if offset < 0 {
		offset = 0
	}

	filters := port.ListParcelFilters{
		Status:              statusPtr,
		OriginOfficeID:      originOfficeID,
		DestinationOfficeID: destinationOfficeID,
		VehicleID:           vehicleID,
		SenderPersonID:      senderPersonID,
		RecipientPersonID:   recipientPersonID,
		FromCreatedAt:       fromPtr,
		ToCreatedAt:         toPtr,
		Limit:               limit,
		Offset:              offset,
	}

	q := strings.TrimSpace(c.Query("q"))
	if q != "" {
		filters.Query = &q
	}

	out, err := h.listUC.Execute(c.Request.Context(), usecase.ListParcelsInput{
		TenantID: tenant,
		Filters:  filters,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	items := make([]dto.CreateParcelResponse, 0, len(out.Items))
	for _, p := range out.Items {
		var registeredAtStr *string
		if p.RegisteredAt != nil {
			s := p.RegisteredAt.UTC().Format(time.RFC3339)
			registeredAtStr = &s
		}
		var boardedAtStr *string
		if p.BoardedAt != nil {
			s := p.BoardedAt.UTC().Format(time.RFC3339)
			boardedAtStr = &s
		}
		var boardedDepartureAtStr *string
		if p.BoardedDepartureAt != nil {
			s := p.BoardedDepartureAt.UTC().Format(time.RFC3339)
			boardedDepartureAtStr = &s
		}
		var departedAtStr *string
		if p.DepartedAt != nil {
			s := p.DepartedAt.UTC().Format(time.RFC3339)
			departedAtStr = &s
		}
		var arrivedAtStr *string
		if p.ArrivedAt != nil {
			s := p.ArrivedAt.UTC().Format(time.RFC3339)
			arrivedAtStr = &s
		}
		var deliveredAtStr *string
		if p.DeliveredAt != nil {
			s := p.DeliveredAt.UTC().Format(time.RFC3339)
			deliveredAtStr = &s
		}

		items = append(items, dto.CreateParcelResponse{
			ID:                  p.ID,
			Status:              string(p.Status),
			ShipmentType:        string(p.ShipmentType),
			OriginOfficeID:      p.OriginOfficeID,
			DestinationOfficeID: p.DestinationOfficeID,
			SenderPersonID:      p.SenderPersonID,
			RecipientPersonID:   p.RecipientPersonID,
			Notes:               p.Notes,
			CreatedAt:           p.CreatedAt.UTC().Format(time.RFC3339),
			RegisteredAt:        registeredAtStr,
			BoardedVehicleID:    p.BoardedVehicleID,
			BoardedTripID:       p.BoardedTripID,
			BoardedDepartureAt:  boardedDepartureAtStr,
			BoardedAt:           boardedAtStr,
			BoardedByUserID:     p.BoardedByUserID,
			DepartedAt:          departedAtStr,
			DepartedByUserID:    p.DepartedByUserID,
			ArrivedAt:           arrivedAtStr,
			ArrivedByUserID:     p.ArrivedByUserID,
			DeliveredAt:         deliveredAtStr,
			DeliveredByUserID:   p.DeliveredByUserID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.ParcelListResponse{
			Items: items,
			Pagination: dto.ParcelListPagination{
				Limit:  limit,
				Offset: offset,
				Count:  out.Count,
			},
		},
	})
}

// GetByID godoc
// @Summary Obtener envío
// @Description Obtiene un envío por id
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.CreateParcelResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id} [get]
func (h *ParcelHandler) GetByID(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")

	p, err := h.getUC.Execute(c.Request.Context(), usecase.GetParcelInput{
		TenantID: strings.TrimSpace(anyToString(tenantID)),
		ParcelID: id,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var registeredAtStr *string
	if p.RegisteredAt != nil {
		s := p.RegisteredAt.UTC().Format(time.RFC3339)
		registeredAtStr = &s
	}
	var boardedAtStr *string
	if p.BoardedAt != nil {
		s := p.BoardedAt.UTC().Format(time.RFC3339)
		boardedAtStr = &s
	}
	var boardedDepartureAtStr *string
	if p.BoardedDepartureAt != nil {
		s := p.BoardedDepartureAt.UTC().Format(time.RFC3339)
		boardedDepartureAtStr = &s
	}
	var arrivedAtStr *string
	if p.ArrivedAt != nil {
		s := p.ArrivedAt.UTC().Format(time.RFC3339)
		arrivedAtStr = &s
	}
	var deliveredAtStr *string
	if p.DeliveredAt != nil {
		s := p.DeliveredAt.UTC().Format(time.RFC3339)
		deliveredAtStr = &s
	}
	var departedAtStr *string
	if p.DepartedAt != nil {
		s := p.DepartedAt.UTC().Format(time.RFC3339)
		departedAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.CreateParcelResponse{
			ID:                  p.ID,
			Status:              string(p.Status),
			ShipmentType:        string(p.ShipmentType),
			OriginOfficeID:      p.OriginOfficeID,
			DestinationOfficeID: p.DestinationOfficeID,
			SenderPersonID:      p.SenderPersonID,
			RecipientPersonID:   p.RecipientPersonID,
			Notes:               p.Notes,
			CreatedAt:           p.CreatedAt.UTC().Format(time.RFC3339),
			RegisteredAt:        registeredAtStr,
			BoardedVehicleID:    p.BoardedVehicleID,
			BoardedTripID:       p.BoardedTripID,
			BoardedDepartureAt:  boardedDepartureAtStr,
			BoardedAt:           boardedAtStr,
			BoardedByUserID:     p.BoardedByUserID,
			ArrivedAt:           arrivedAtStr,
			ArrivedByUserID:     p.ArrivedByUserID,
			DeliveredAt:         deliveredAtStr,
			DeliveredByUserID:   p.DeliveredByUserID,
			DepartedAt:          departedAtStr,
			DepartedByUserID:    p.DepartedByUserID,
		},
	})
}

// Register godoc
// @Summary Registrar envío
// @Description Cambia estado a registrado
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.CreateParcelResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/register [post]
func (h *ParcelHandler) Register(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	p, err := h.registerUC.Execute(c.Request.Context(), usecase.RegisterParcelInput{
		TenantID: strings.TrimSpace(anyToString(tenantID)),
		UserID:   strings.TrimSpace(anyToString(userID)),
		UserName: strings.TrimSpace(anyToString(userName)),
		ParcelID: id,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var registeredAtStr *string
	if p.RegisteredAt != nil {
		s := p.RegisteredAt.UTC().Format(time.RFC3339)
		registeredAtStr = &s
	}
	var boardedAtStr *string
	if p.BoardedAt != nil {
		s := p.BoardedAt.UTC().Format(time.RFC3339)
		boardedAtStr = &s
	}
	var boardedDepartureAtStr *string
	if p.BoardedDepartureAt != nil {
		s := p.BoardedDepartureAt.UTC().Format(time.RFC3339)
		boardedDepartureAtStr = &s
	}
	var arrivedAtStr *string
	if p.ArrivedAt != nil {
		s := p.ArrivedAt.UTC().Format(time.RFC3339)
		arrivedAtStr = &s
	}
	var deliveredAtStr *string
	if p.DeliveredAt != nil {
		s := p.DeliveredAt.UTC().Format(time.RFC3339)
		deliveredAtStr = &s
	}
	var departedAtStr *string
	if p.DepartedAt != nil {
		s := p.DepartedAt.UTC().Format(time.RFC3339)
		departedAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.CreateParcelResponse{
			ID:                  p.ID,
			Status:              string(p.Status),
			ShipmentType:        string(p.ShipmentType),
			OriginOfficeID:      p.OriginOfficeID,
			DestinationOfficeID: p.DestinationOfficeID,
			SenderPersonID:      p.SenderPersonID,
			RecipientPersonID:   p.RecipientPersonID,
			Notes:               p.Notes,
			CreatedAt:           p.CreatedAt.UTC().Format(time.RFC3339),
			RegisteredAt:        registeredAtStr,
			BoardedVehicleID:    p.BoardedVehicleID,
			BoardedTripID:       p.BoardedTripID,
			BoardedDepartureAt:  boardedDepartureAtStr,
			BoardedAt:           boardedAtStr,
			BoardedByUserID:     p.BoardedByUserID,
			ArrivedAt:           arrivedAtStr,
			ArrivedByUserID:     p.ArrivedByUserID,
			DeliveredAt:         deliveredAtStr,
			DeliveredByUserID:   p.DeliveredByUserID,
			DepartedAt:          departedAtStr,
			DepartedByUserID:    p.DepartedByUserID,
		},
	})
}

// Board godoc
// @Summary Embarcar envío
// @Description Cambia estado a embarcado
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.CreateParcelResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/board [post]
func (h *ParcelHandler) Board(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	var req dto.BoardParcelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	vehicleUUID, err := uuid.Parse(strings.TrimSpace(req.VehicleID))
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "vehicle_id inválido", map[string]any{"field": "vehicle_id"}))
		return
	}

	// Validación opcional de consistencia de origin_office_id (sin integrar LOCATION)
	if req.OriginOfficeID != nil {
		if _, err := uuid.Parse(strings.TrimSpace(*req.OriginOfficeID)); err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "origin_office_id inválido", map[string]any{"field": "origin_office_id"}))
			return
		}
	}

	var tripUUID *uuid.UUID
	if req.TripID != nil && strings.TrimSpace(*req.TripID) != "" {
		t, err := uuid.Parse(strings.TrimSpace(*req.TripID))
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "trip_id inválido", map[string]any{"field": "trip_id"}))
			return
		}
		tripUUID = &t
	}

	var departureAt *time.Time
	if req.DepartureAt != nil && strings.TrimSpace(*req.DepartureAt) != "" {
		tm, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.DepartureAt))
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "departure_at inválido", map[string]any{"field": "departure_at"}))
			return
		}
		ut := tm.UTC()
		departureAt = &ut
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	p, err := h.boardUC.Execute(c.Request.Context(), usecase.BoardParcelInput{
		TenantID:    strings.TrimSpace(anyToString(tenantID)),
		UserID:      strings.TrimSpace(anyToString(userID)),
		UserName:    strings.TrimSpace(anyToString(userName)),
		ParcelID:    id,
		VehicleID:   vehicleUUID,
		TripID:      tripUUID,
		DepartureAt: departureAt,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var registeredAtStr *string
	if p.RegisteredAt != nil {
		s := p.RegisteredAt.UTC().Format(time.RFC3339)
		registeredAtStr = &s
	}
	var boardedAtStr *string
	if p.BoardedAt != nil {
		s := p.BoardedAt.UTC().Format(time.RFC3339)
		boardedAtStr = &s
	}
	var boardedDepartureAtStr *string
	if p.BoardedDepartureAt != nil {
		s := p.BoardedDepartureAt.UTC().Format(time.RFC3339)
		boardedDepartureAtStr = &s
	}
	var arrivedAtStr *string
	if p.ArrivedAt != nil {
		s := p.ArrivedAt.UTC().Format(time.RFC3339)
		arrivedAtStr = &s
	}
	var deliveredAtStr *string
	if p.DeliveredAt != nil {
		s := p.DeliveredAt.UTC().Format(time.RFC3339)
		deliveredAtStr = &s
	}
	var departedAtStr *string
	if p.DepartedAt != nil {
		s := p.DepartedAt.UTC().Format(time.RFC3339)
		departedAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.CreateParcelResponse{
			ID:                  p.ID,
			Status:              string(p.Status),
			ShipmentType:        string(p.ShipmentType),
			OriginOfficeID:      p.OriginOfficeID,
			DestinationOfficeID: p.DestinationOfficeID,
			SenderPersonID:      p.SenderPersonID,
			RecipientPersonID:   p.RecipientPersonID,
			Notes:               p.Notes,
			CreatedAt:           p.CreatedAt.UTC().Format(time.RFC3339),
			RegisteredAt:        registeredAtStr,
			BoardedVehicleID:    p.BoardedVehicleID,
			BoardedTripID:       p.BoardedTripID,
			BoardedDepartureAt:  boardedDepartureAtStr,
			BoardedAt:           boardedAtStr,
			BoardedByUserID:     p.BoardedByUserID,
			DepartedAt:          departedAtStr,
			DepartedByUserID:    p.DepartedByUserID,
			ArrivedAt:           arrivedAtStr,
			ArrivedByUserID:     p.ArrivedByUserID,
			DeliveredAt:         deliveredAtStr,
			DeliveredByUserID:   p.DeliveredByUserID,
		},
	})
}

// Depart godoc
// @Summary Salida de envío
// @Description Cambia estado a en ruta
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.CreateParcelResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/depart [post]
func (h *ParcelHandler) Depart(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	var req dto.DepartParcelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	var vehicleUUID *uuid.UUID
	if req.VehicleID != nil && strings.TrimSpace(*req.VehicleID) != "" {
		v, err := uuid.Parse(strings.TrimSpace(*req.VehicleID))
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "vehicle_id inválido", map[string]any{"field": "vehicle_id"}))
			return
		}
		vehicleUUID = &v
	}

	var departedAt *time.Time
	if req.DepartedAt != nil && strings.TrimSpace(*req.DepartedAt) != "" {
		tm, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.DepartedAt))
		if err != nil {
			_ = c.Error(apperror.NewBadRequest("validation_error", "departed_at inválido", map[string]any{"field": "departed_at"}))
			return
		}
		ut := tm.UTC()
		departedAt = &ut
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	p, err := h.departUC.Execute(c.Request.Context(), usecase.DepartParcelInput{
		TenantID:          tenant,
		UserID:            strings.TrimSpace(anyToString(userID)),
		UserName:          strings.TrimSpace(anyToString(userName)),
		ParcelID:          id,
		DepartureOfficeID: req.DepartureOfficeID,
		VehicleID:         vehicleUUID,
		DepartedAt:        departedAt,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var registeredAtStr *string
	if p.RegisteredAt != nil {
		s := p.RegisteredAt.UTC().Format(time.RFC3339)
		registeredAtStr = &s
	}
	var departedAtStr *string
	if p.DepartedAt != nil {
		s := p.DepartedAt.UTC().Format(time.RFC3339)
		departedAtStr = &s
	}
	var boardedAtStr *string
	if p.BoardedAt != nil {
		s := p.BoardedAt.UTC().Format(time.RFC3339)
		boardedAtStr = &s
	}
	var boardedDepartureAtStr *string
	if p.BoardedDepartureAt != nil {
		s := p.BoardedDepartureAt.UTC().Format(time.RFC3339)
		boardedDepartureAtStr = &s
	}
	var arrivedAtStr *string
	if p.ArrivedAt != nil {
		s := p.ArrivedAt.UTC().Format(time.RFC3339)
		arrivedAtStr = &s
	}
	var deliveredAtStr *string
	if p.DeliveredAt != nil {
		s := p.DeliveredAt.UTC().Format(time.RFC3339)
		deliveredAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.CreateParcelResponse{
			ID:                  p.ID,
			Status:              string(p.Status),
			ShipmentType:        string(p.ShipmentType),
			OriginOfficeID:      p.OriginOfficeID,
			DestinationOfficeID: p.DestinationOfficeID,
			SenderPersonID:      p.SenderPersonID,
			RecipientPersonID:   p.RecipientPersonID,
			Notes:               p.Notes,
			CreatedAt:           p.CreatedAt.UTC().Format(time.RFC3339),
			RegisteredAt:        registeredAtStr,
			BoardedVehicleID:    p.BoardedVehicleID,
			BoardedTripID:       p.BoardedTripID,
			BoardedDepartureAt:  boardedDepartureAtStr,
			BoardedAt:           boardedAtStr,
			BoardedByUserID:     p.BoardedByUserID,
			DepartedAt:          departedAtStr,
			DepartedByUserID:    p.DepartedByUserID,
			ArrivedAt:           arrivedAtStr,
			ArrivedByUserID:     p.ArrivedByUserID,
			DeliveredAt:         deliveredAtStr,
			DeliveredByUserID:   p.DeliveredByUserID,
		},
	})
}

// Arrive godoc
// @Summary Llegada a destino
// @Description Cambia estado a en oficina destino
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.CreateParcelResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/arrive [post]
func (h *ParcelHandler) Arrive(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	var req dto.ArriveParcelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	p, err := h.arriveUC.Execute(c.Request.Context(), usecase.ArriveParcelInput{
		TenantID:            tenant,
		UserID:              strings.TrimSpace(anyToString(userID)),
		UserName:            strings.TrimSpace(anyToString(userName)),
		ParcelID:            id,
		DestinationOfficeID: req.DestinationOfficeID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var registeredAtStr *string
	if p.RegisteredAt != nil {
		s := p.RegisteredAt.UTC().Format(time.RFC3339)
		registeredAtStr = &s
	}
	var boardedAtStr *string
	if p.BoardedAt != nil {
		s := p.BoardedAt.UTC().Format(time.RFC3339)
		boardedAtStr = &s
	}
	var boardedDepartureAtStr *string
	if p.BoardedDepartureAt != nil {
		s := p.BoardedDepartureAt.UTC().Format(time.RFC3339)
		boardedDepartureAtStr = &s
	}
	var arrivedAtStr *string
	if p.ArrivedAt != nil {
		s := p.ArrivedAt.UTC().Format(time.RFC3339)
		arrivedAtStr = &s
	}
	var deliveredAtStr *string
	if p.DeliveredAt != nil {
		s := p.DeliveredAt.UTC().Format(time.RFC3339)
		deliveredAtStr = &s
	}
	var departedAtStr *string
	if p.DepartedAt != nil {
		s := p.DepartedAt.UTC().Format(time.RFC3339)
		departedAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.CreateParcelResponse{
			ID:                  p.ID,
			Status:              string(p.Status),
			ShipmentType:        string(p.ShipmentType),
			OriginOfficeID:      p.OriginOfficeID,
			DestinationOfficeID: p.DestinationOfficeID,
			SenderPersonID:      p.SenderPersonID,
			RecipientPersonID:   p.RecipientPersonID,
			Notes:               p.Notes,
			CreatedAt:           p.CreatedAt.UTC().Format(time.RFC3339),
			RegisteredAt:        registeredAtStr,
			BoardedVehicleID:    p.BoardedVehicleID,
			BoardedTripID:       p.BoardedTripID,
			BoardedDepartureAt:  boardedDepartureAtStr,
			BoardedAt:           boardedAtStr,
			BoardedByUserID:     p.BoardedByUserID,
			DepartedAt:          departedAtStr,
			DepartedByUserID:    p.DepartedByUserID,
			ArrivedAt:           arrivedAtStr,
			ArrivedByUserID:     p.ArrivedByUserID,
			DeliveredAt:         deliveredAtStr,
			DeliveredByUserID:   p.DeliveredByUserID,
		},
	})
}

// Deliver godoc
// @Summary Entregar envío
// @Description Cambia estado a entregado
// @Tags Parcels
// @Produce json
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Param id path string true "UUID" Format(uuid)
// @Success 200 {object} handler.CreateParcelResponseEnvelope
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 409 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /parcels/{id}/deliver [post]
func (h *ParcelHandler) Deliver(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "id inválido", map[string]any{"field": "id"}))
		return
	}

	var req dto.DeliverParcelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewBadRequest("validation_error", "payload inválido", map[string]any{"error": err.Error()}))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	tenant := strings.TrimSpace(anyToString(tenantID))
	if tenant == "" {
		_ = c.Error(apperror.NewUnauthorized("unauthorized", "credenciales inválidas", nil))
		return
	}

	p, err := h.deliverUC.Execute(c.Request.Context(), usecase.DeliverParcelInput{
		TenantID:   tenant,
		UserID:     strings.TrimSpace(anyToString(userID)),
		UserName:   strings.TrimSpace(anyToString(userName)),
		ParcelID:   id,
		PackageKey: req.PackageKey,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var registeredAtStr *string
	if p.RegisteredAt != nil {
		s := p.RegisteredAt.UTC().Format(time.RFC3339)
		registeredAtStr = &s
	}
	var boardedAtStr *string
	if p.BoardedAt != nil {
		s := p.BoardedAt.UTC().Format(time.RFC3339)
		boardedAtStr = &s
	}
	var boardedDepartureAtStr *string
	if p.BoardedDepartureAt != nil {
		s := p.BoardedDepartureAt.UTC().Format(time.RFC3339)
		boardedDepartureAtStr = &s
	}
	var arrivedAtStr *string
	if p.ArrivedAt != nil {
		s := p.ArrivedAt.UTC().Format(time.RFC3339)
		arrivedAtStr = &s
	}
	var deliveredAtStr *string
	if p.DeliveredAt != nil {
		s := p.DeliveredAt.UTC().Format(time.RFC3339)
		deliveredAtStr = &s
	}
	var departedAtStr *string
	if p.DepartedAt != nil {
		s := p.DepartedAt.UTC().Format(time.RFC3339)
		departedAtStr = &s
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.CreateParcelResponse{
			ID:                  p.ID,
			Status:              string(p.Status),
			ShipmentType:        string(p.ShipmentType),
			OriginOfficeID:      p.OriginOfficeID,
			DestinationOfficeID: p.DestinationOfficeID,
			SenderPersonID:      p.SenderPersonID,
			RecipientPersonID:   p.RecipientPersonID,
			Notes:               p.Notes,
			CreatedAt:           p.CreatedAt.UTC().Format(time.RFC3339),
			RegisteredAt:        registeredAtStr,
			BoardedVehicleID:    p.BoardedVehicleID,
			BoardedTripID:       p.BoardedTripID,
			BoardedDepartureAt:  boardedDepartureAtStr,
			BoardedAt:           boardedAtStr,
			BoardedByUserID:     p.BoardedByUserID,
			DepartedAt:          departedAtStr,
			DepartedByUserID:    p.DepartedByUserID,
			ArrivedAt:           arrivedAtStr,
			ArrivedByUserID:     p.ArrivedByUserID,
			DeliveredAt:         deliveredAtStr,
			DeliveredByUserID:   p.DeliveredByUserID,
		},
	})
}

func anyToString(v any) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return ""
}
