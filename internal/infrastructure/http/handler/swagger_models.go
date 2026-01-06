package handler

import "ms-parcel-core/internal/pkg/util/apperror"

type ErrorResponse struct {
	Success bool              `json:"success" example:"false"`
	Error   apperror.AppError `json:"error"`
}

type AnyDataEnvelope struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data"`
}

type AnyListEnvelope struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data"`
}

// System

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

type HealthResponseEnvelope struct {
	Data HealthResponse `json:"data"`
}
