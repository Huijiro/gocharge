package gocharge

import (
	"context"
	"net/http"
)

// ErrorCode is a typed string representing a specific error category.
//
// Error codes provide semantic information about what went wrong in a way
// that's easy for clients to parse and handle programmatically.
// Each ErrorCode has an associated HTTP status code.
//
// Example:
//
//	code := gocharge.ErrNotFound
//	appErr := gocharge.NewError(code, "User not found", http.StatusNotFound)
type ErrorCode string

// Standard error codes for common HTTP scenarios.
//
// These codes are sent to clients in error responses and can be used
// to make error handling decisions. Each code maps to a standard HTTP status.
const (
	// ErrBadRequest (400) - Invalid request format or parameters
	ErrBadRequest ErrorCode = "BAD_REQUEST"
	// ErrUnauthorized (401) - Request requires authentication
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrForbidden (403) - Authenticated but not authorized
	ErrForbidden ErrorCode = "FORBIDDEN"
	// ErrNotFound (404) - Resource does not exist
	ErrNotFound ErrorCode = "NOT_FOUND"
	// ErrConflict (409) - Resource conflict (e.g., duplicate)
	ErrConflict ErrorCode = "CONFLICT"
	// ErrValidation (400) - Validation failure
	ErrValidation ErrorCode = "VALIDATION_ERROR"
	// ErrInternalServerError (500) - Unexpected server error
	ErrInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	// ErrServiceUnavailable (503) - Service temporarily unavailable
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	// ErrUnprocessableEntity (422) - Valid syntax but semantic error
	ErrUnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY"
	// ErrTooManyRequests (429) - Rate limit exceeded
	ErrTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
)

// Errorable defines the interface for typed errors in gocharge.
//
// All errors returned from handlers should either implement this interface
// or be handled by a custom ErrorEncoder. The framework uses this interface
// to determine HTTP status codes and error codes for responses.
//
// Example:
//
//	type MyError struct {
//	    code   gocharge.ErrorCode
//	    msg    string
//	    status int
//	}
//	func (e MyError) Error() string { return e.msg }
//	func (e MyError) Code() gocharge.ErrorCode { return e.code }
//	func (e MyError) StatusCode() int { return e.status }
type Errorable interface {
	error
	// Code returns the error code for this error
	Code() ErrorCode
	// StatusCode returns the HTTP status code for this error
	StatusCode() int
}

// AppError is a simple, built-in implementation of the Errorable interface.
//
// It stores an error code, message, and HTTP status code, providing a lightweight
// way to return typed errors from handlers without defining custom error types.
type AppError struct {
	code       ErrorCode
	message    string
	statusCode int
}

// NewError creates a new AppError with explicit code, message, and status.
//
// This is the low-level constructor. For convenience, use NewAppError which
// automatically maps error codes to standard HTTP status codes.
//
// Example:
//
//	err := gocharge.NewError(
//	    gocharge.ErrNotFound,
//	    "User not found",
//	    http.StatusNotFound,
//	)
func NewError(code ErrorCode, message string, statusCode int) AppError {
	return AppError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

// Error implements the error interface, returning the message.
func (e AppError) Error() string {
	return e.message
}

// Code returns the ErrorCode for this error.
func (e AppError) Code() ErrorCode {
	return e.code
}

// StatusCode returns the HTTP status code for this error.
func (e AppError) StatusCode() int {
	return e.statusCode
}

// ErrorResponse is the JSON structure sent to clients when an error occurs.
//
// It contains:
//   - Code: Semantic error code (e.g., "NOT_FOUND")
//   - Message: Human-readable error message
//
// Example JSON:
//
//	{
//	  "code": "NOT_FOUND",
//	  "message": "User not found"
//	}
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorEncoder is responsible for converting errors to ErrorResponse JSON.
//
// Implementers can customize how errors are presented to clients.
// The encoder receives the request context (for logging context) and the error,
// and returns both the ErrorResponse and HTTP status code.
//
// This allows applications to:
//   - Map custom error types to error codes
//   - Include additional fields in responses
//   - Apply different encoding based on context
//   - Implement logging or metrics
//
// Example:
//
//	type MyErrorEncoder struct{}
//	func (e *MyErrorEncoder) Encode(ctx context.Context, err error) (ErrorResponse, int) {
//	    if appErr, ok := err.(gocharge.Errorable); ok {
//	        return ErrorResponse{
//	            Code:    string(appErr.Code()),
//	            Message: appErr.Error(),
//	        }, appErr.StatusCode()
//	    }
//	    return ErrorResponse{
//	        Code:    "UNKNOWN_ERROR",
//	        Message: "An unexpected error occurred",
//	    }, http.StatusInternalServerError
//	}
type ErrorEncoder interface {
	// Encode converts an error to an ErrorResponse and HTTP status code.
	// The context may contain request-scoped values useful for encoding.
	Encode(ctx context.Context, err error) (ErrorResponse, int)
}

// DefaultErrorEncoder converts errors to standard ErrorResponse format.
//
// If the error implements Errorable, it uses the error's Code() and StatusCode().
// For non-Errorable errors, it returns a generic 500 error.
//
// This provides a reasonable default for most applications, but can be replaced
// with a custom encoder via Server.SetErrorEncoder().
type DefaultErrorEncoder struct{}

// Encode converts an error to ErrorResponse and status code.
//
// For Errorable errors, extracts the code and status.
// For other errors, returns a generic 500 response.
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

// NewAppError creates an AppError with automatic HTTP status code mapping.
//
// This is a convenience function that maps ErrorCode values to standard HTTP status codes.
// For example, ErrNotFound maps to 404, ErrValidation maps to 400, etc.
//
// This is the recommended way to create AppError instances in handler code.
//
// Example:
//
//	if user == nil {
//	    return gocharge.NewAppError(gocharge.ErrNotFound, "User not found")
//	}
//	if !isValid(req) {
//	    return gocharge.NewAppError(gocharge.ErrValidation, "Invalid input")
//	}
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
