package gocharge

import (
	"encoding/json"
	"net/http"
)

// Response[T] wraps an http.ResponseWriter with a typed response body field.
//
// Type parameter T is the type of data to be sent in the response JSON.
// The embedded http.ResponseWriter provides access to HTTP-specific methods
// like setting headers.
//
// Methods support chaining for fluent API design:
//
//	_, err := w.Status(gocharge.StatusCreated).JSON(newUser)
//	return err
//
// Fields:
//   - StatusCode: The HTTP status code to send (default 200)
//   - Data: The response body data (set by JSON())
//   - http.ResponseWriter: The underlying response writer
type Response[T any] struct {
	StatusCode int
	Data       T
	http.ResponseWriter
}

// StatusCode is a typed HTTP status code.
//
// Using a distinct type allows the package to provide convenient constants
// for all standard HTTP status codes. Functions requiring http.StatusCode
// from the standard library can use int(gocharge.StatusCreated), etc.
type StatusCode int

const (
	// 1xx Informational
	StatusContinue           StatusCode = 100
	StatusSwitchingProtocols StatusCode = 101
	StatusProcessing         StatusCode = 102
	StatusEarlyHints         StatusCode = 103

	// 2xx Success
	StatusOK                   StatusCode = 200
	StatusCreated              StatusCode = 201
	StatusAccepted             StatusCode = 202
	StatusNonAuthoritativeInfo StatusCode = 203
	StatusNoContent            StatusCode = 204
	StatusResetContent         StatusCode = 205
	StatusPartialContent       StatusCode = 206
	StatusMultiStatus          StatusCode = 207
	StatusAlreadyReported      StatusCode = 208
	StatusIMUsed               StatusCode = 226

	// 3xx Redirection
	StatusMultipleChoices   StatusCode = 300
	StatusMovedPermanently  StatusCode = 301
	StatusFound             StatusCode = 302
	StatusSeeOther          StatusCode = 303
	StatusNotModified       StatusCode = 304
	StatusUseProxy          StatusCode = 305
	StatusTemporaryRedirect StatusCode = 307
	StatusPermanentRedirect StatusCode = 308

	// 4xx Client Error
	StatusBadRequest                   StatusCode = 400
	StatusUnauthorized                 StatusCode = 401
	StatusPaymentRequired              StatusCode = 402
	StatusForbidden                    StatusCode = 403
	StatusNotFound                     StatusCode = 404
	StatusMethodNotAllowed             StatusCode = 405
	StatusNotAcceptable                StatusCode = 406
	StatusProxyAuthRequired            StatusCode = 407
	StatusRequestTimeout               StatusCode = 408
	StatusConflict                     StatusCode = 409
	StatusGone                         StatusCode = 410
	StatusLengthRequired               StatusCode = 411
	StatusPreconditionFailed           StatusCode = 412
	StatusRequestEntityTooLarge        StatusCode = 413
	StatusRequestURITooLong            StatusCode = 414
	StatusUnsupportedMediaType         StatusCode = 415
	StatusRequestedRangeNotSatisfiable StatusCode = 416
	StatusExpectationFailed            StatusCode = 417
	StatusTeapot                       StatusCode = 418
	StatusMisdirectedRequest           StatusCode = 421
	StatusUnprocessableEntity          StatusCode = 422
	StatusLocked                       StatusCode = 423
	StatusFailedDependency             StatusCode = 424
	StatusTooEarly                     StatusCode = 425
	StatusUpgradeRequired              StatusCode = 426
	StatusPreconditionRequired         StatusCode = 428
	StatusTooManyRequests              StatusCode = 429
	StatusRequestHeaderFieldsTooLarge  StatusCode = 431
	StatusUnavailableForLegalReasons   StatusCode = 451

	// 5xx Server Error
	StatusInternalServerError           StatusCode = 500
	StatusNotImplemented                StatusCode = 501
	StatusBadGateway                    StatusCode = 502
	StatusServiceUnavailable            StatusCode = 503
	StatusGatewayTimeout                StatusCode = 504
	StatusHTTPVersionNotSupported       StatusCode = 505
	StatusVariantAlsoNegotiates         StatusCode = 506
	StatusInsufficientStorage           StatusCode = 507
	StatusLoopDetected                  StatusCode = 508
	StatusNotExtended                   StatusCode = 510
	StatusNetworkAuthenticationRequired StatusCode = 511
)

// Status sets the HTTP status code and returns the receiver for method chaining.
//
// If Status is not called, the default status code 200 OK is used.
// Status must be called before JSON() or Error() to set a non-200 status code.
//
// Returns the receiver to allow chaining:
//
//	_, err := w.Status(gocharge.StatusCreated).JSON(user)
//	_, err := w.Status(gocharge.StatusBadRequest).Error("INVALID", "Invalid input")
func (r *Response[T]) Status(statusCode StatusCode) *Response[T] {
	r.StatusCode = int(statusCode)
	return r
}

// JSON writes the typed response data as JSON.
//
// The data is encoded as JSON with Content-Type: application/json header.
// The previously set status code (or default 200) is sent.
// Returns the receiver and any error from JSON encoding.
//
// Once JSON is called, the HTTP response is sent and cannot be modified.
//
// Common patterns:
//
//	// Success response
//	_, err := w.JSON(user)
//	return err
//
//	// With custom status
//	_, err := w.Status(gocharge.StatusCreated).JSON(newUser)
//	return err
//
//	// With multiple headers (set before calling JSON)
//	w.Header().Set("X-Total-Count", "100")
//	_, err := w.JSON(items)
//	return err
func (r *Response[T]) JSON(data T) (*Response[T], error) {
	r.Data = data
	r.Header().Set("Content-Type", "application/json")

	if r.StatusCode != 0 {
		r.WriteHeader(r.StatusCode)
	} else {
		r.WriteHeader(http.StatusOK)
	}

	err := json.NewEncoder(r).Encode(data)
	return r, err
}

// Error writes an error response with the given code and message.
//
// This is useful for handlers that want to send error responses without
// returning from the handler function. The error is encoded as ErrorResponse JSON.
//
// If Status was called, that status code is used. Otherwise defaults to 500.
// Once Error is called, the HTTP response is sent and cannot be modified.
//
// Note: Most handlers should return an error instead of calling Error(),
// which allows the server's ErrorEncoder to handle error responses consistently.
//
// Example:
//
//	if !hasPermission {
//	    _, err := w.Status(gocharge.StatusForbidden).Error("PERMISSION_DENIED", "Access denied")
//	    return err
//	}
func (r *Response[T]) Error(code string, message string) (*Response[T], error) {
	r.Header().Set("Content-Type", "application/json")

	if r.StatusCode == 0 {
		r.WriteHeader(http.StatusInternalServerError)
	} else {
		r.WriteHeader(r.StatusCode)
	}

	errResp := ErrorResponse{
		Code:    code,
		Message: message,
	}

	err := json.NewEncoder(r).Encode(errResp)
	return r, err
}
