package gocharge

import (
	"context"
	"net/http"
)

// ErrorCode represents a typed error code for API errors
type ErrorCode string

// Standard error codes
const (
	ErrBadRequest          ErrorCode = "BAD_REQUEST"
	ErrUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrForbidden           ErrorCode = "FORBIDDEN"
	ErrNotFound            ErrorCode = "NOT_FOUND"
	ErrConflict            ErrorCode = "CONFLICT"
	ErrValidation          ErrorCode = "VALIDATION_ERROR"
	ErrInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"
	ErrUnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY"
	ErrTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"
)

// Errorable interface defines what an error must provide for the framework
type Errorable interface {
	error
	Code() ErrorCode
	StatusCode() int
}

// AppError is the default error implementation
type AppError struct {
	code       ErrorCode
	message    string
	statusCode int
}

// NewError creates a new AppError with the given code, message, and HTTP status code
func NewError(code ErrorCode, message string, statusCode int) AppError {
	return AppError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

// Error implements the error interface
func (e AppError) Error() string {
	return e.message
}

// Code returns the error code
func (e AppError) Code() ErrorCode {
	return e.code
}

// StatusCode returns the HTTP status code
func (e AppError) StatusCode() int {
	return e.statusCode
}

// ErrorResponse is the JSON response sent to clients on error
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorEncoder is implemented by types that can encode errors for API responses
type ErrorEncoder interface {
	Encode(ctx context.Context, err error) (ErrorResponse, int)
}

// DefaultErrorEncoder is the default implementation of ErrorEncoder
type DefaultErrorEncoder struct{}

// Encode converts an error to an ErrorResponse and HTTP status code
func (enc *DefaultErrorEncoder) Encode(ctx context.Context, err error) (ErrorResponse, int) {
	// Check if error implements Errorable
	if appErr, ok := err.(Errorable); ok {
		return ErrorResponse{
			Code:    string(appErr.Code()),
			Message: appErr.Error(),
		}, appErr.StatusCode()
	}

	// For generic errors, return 500
	return ErrorResponse{
		Code:    string(ErrInternalServerError),
		Message: "An internal server error occurred",
	}, http.StatusInternalServerError
}

// NewAppError creates a new AppError with standard HTTP status codes based on error code
// This is a convenience function for common error patterns
func NewAppError(code ErrorCode, message string) AppError {
	statusCode := http.StatusInternalServerError

	switch code {
	case ErrBadRequest, ErrValidation:
		statusCode = http.StatusBadRequest
	case ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case ErrForbidden:
		statusCode = http.StatusForbidden
	case ErrNotFound:
		statusCode = http.StatusNotFound
	case ErrConflict:
		statusCode = http.StatusConflict
	case ErrUnprocessableEntity:
		statusCode = http.StatusUnprocessableEntity
	case ErrTooManyRequests:
		statusCode = http.StatusTooManyRequests
	case ErrServiceUnavailable:
		statusCode = http.StatusServiceUnavailable
	}

	return NewError(code, message, statusCode)
}
