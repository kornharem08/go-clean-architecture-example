package pkg

import "errors"

// Common error messages
var (
	ErrInvalidInput = errors.New("invalid input provided")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInternal     = errors.New("internal server error")
)

// AppError is a custom error type for the application
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewAppError creates a new AppError
func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}
