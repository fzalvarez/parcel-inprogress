package apperror

import "net/http"

// AppError represents an error with an HTTP status code.
type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Status  int         `json:"-"`
}

// Error makes AppError implement the built-in error interface.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

// New creates a new AppError instance.
func New(code string, message string, details interface{}, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Status:  status,
	}
}

// NewBadRequest creates a new AppError with HTTP status 400 (Bad Request).
func NewBadRequest(code string, message string, details any) *AppError {
	return New(code, message, details, http.StatusBadRequest)
}

// NewUnauthorized creates a new AppError with HTTP status 401 (Unauthorized).
func NewUnauthorized(code string, message string, details any) *AppError {
	return New(code, message, details, http.StatusUnauthorized)
}

// NewInternal creates a new AppError with HTTP status 500 (Internal Server Error).
func NewInternal(code string, message string, details any) *AppError {
	return New(code, message, details, http.StatusInternalServerError)
}
