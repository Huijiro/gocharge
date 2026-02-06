package gocharge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// PathParamConverter[V] converts a string path parameter to type V.
//
// Implementations should return a detailed error message if conversion fails,
// following the pattern: "Invalid {paramName} format: expected {type}, got '{value}'"
//
// Example:
//
//	type IntConverter struct{}
//	func (c IntConverter) Convert(paramName, value string) (int, error) {
//	    parsed, err := strconv.Atoi(value)
//	    if err != nil {
//	        return 0, fmt.Errorf("Invalid %s format: expected integer, got '%s'", paramName, value)
//	    }
//	    return parsed, nil
//	}
type PathParamConverter[V any] interface {
	Convert(paramName, value string) (V, error)
}

// StringConverter converts path parameters to strings (passthrough).
type StringConverter struct{}

func (c StringConverter) Convert(_ string, value string) (string, error) {
	return value, nil
}

// IntConverter converts path parameters to int with detailed error messages.
type IntConverter struct{}

func (c IntConverter) Convert(paramName, value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: expected integer, got '%s'", paramName, value)
	}
	return parsed, nil
}

// Int64Converter converts path parameters to int64 with detailed error messages.
type Int64Converter struct{}

func (c Int64Converter) Convert(paramName, value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: expected integer, got '%s'", paramName, value)
	}
	return parsed, nil
}

// Request[T] wraps an http.Request with a typed body data field.
//
// Type parameter T is the expected type of the request body (JSON).
// The embedded http.Request provides access to HTTP-specific data like
// headers, query parameters, and the raw request.
//
// Fields:
//   - Data: The decoded request body
//   - http.Request: The underlying HTTP request
//
// Path parameters are accessed via the PathParam() method:
//   - r.PathParam().String("id") - extract string parameter
//   - r.PathParam().Int("id") - extract integer parameter
//   - r.PathParam().Int64("id") - extract int64 parameter
//   - r.PathParam().UUID("id") - extract and validate UUID parameter
//
// Example:
//
//	type CreateUserRequest struct {
//	    Name string `json:"name"`
//	}
//
//	func handler(ctx context.Context, w Response[Response], r Request[CreateUserRequest]) error {
//	    body, err := r.JSON()
//	    if err != nil {
//	        return gocharge.NewAppError(gocharge.ErrBadRequest, "Invalid JSON")
//	    }
//	    // body is now typed as *CreateUserRequest
//	    return nil
//	}
type Request[T any] struct {
	Data T
	http.Request
}

// JSON decodes the request body as JSON into the Data field.
//
// The decoded value is both stored in r.Data and returned.
// If decoding fails (e.g., invalid JSON), an error is returned and Data
// remains as its zero value.
//
// Returns:
//   - A pointer to the decoded data (r.Data)
//   - An error if JSON decoding fails
//
// Example:
//
//	body, err := r.JSON()
//	if err != nil {
//	    return gocharge.NewAppError(gocharge.ErrBadRequest, "Invalid JSON")
//	}
//	// body is a pointer to the typed data
func (r *Request[T]) JSON() (*T, error) {
	err := json.NewDecoder(r.Body).Decode(&r.Data)

	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}

// Context returns the request's context.
//
// This is a convenience method to access the request context without
// needing to call r.Request.Context(). The context can be used to:
//   - Access request-scoped values set by middleware
//   - Check for request deadlines/timeouts
//   - Pass context to business logic
//
// Example:
//
//	ctx := r.Context()
//	user, hasUser := gocharge.UserFromContext(ctx)
//	logger := gocharge.GetLogger()
//	logger.Info(ctx, "processing request")
func (r *Request[T]) Context() context.Context {
	return r.Request.Context()
}

// PathParamHelper provides type-safe path parameter extraction.
type PathParamHelper struct {
	request *http.Request
}

// pathParamHelper creates a helper for extracting typed path parameters.
func (r *Request[T]) pathParamHelper() PathParamHelper {
	return PathParamHelper{request: &r.Request}
}

// String extracts a string path parameter using Go 1.22+ PathValue.
//
// Returns:
//   - The path parameter value as string
//   - An error if the parameter is empty (not found in route)
//
// Example:
//
//	userID, err := r.PathParam().String("id")
func (h PathParamHelper) String(key string) (string, error) {
	value := h.request.PathValue(key)
	if value == "" {
		return "", fmt.Errorf("missing path parameter: %s", key)
	}
	return value, nil
}

// Int extracts an integer path parameter.
//
// Returns:
//   - The path parameter value as int
//   - An error with detailed message if conversion fails: "invalid {paramName} format: expected integer, got '{value}'"
//
// Example:
//
//	itemID, err := r.PathParam().Int("id")
func (h PathParamHelper) Int(key string) (int, error) {
	value := h.request.PathValue(key)
	if value == "" {
		return 0, fmt.Errorf("missing path parameter: %s", key)
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: expected integer, got '%s'", key, value)
	}
	return parsed, nil
}

// Int64 extracts an int64 path parameter.
//
// Returns:
//   - The path parameter value as int64
//   - An error with detailed message if conversion fails: "invalid {paramName} format: expected integer, got '{value}'"
//
// Example:
//
//	itemID, err := r.PathParam().Int64("id")
func (h PathParamHelper) Int64(key string) (int64, error) {
	value := h.request.PathValue(key)
	if value == "" {
		return 0, fmt.Errorf("missing path parameter: %s", key)
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: expected integer, got '%s'", key, value)
	}
	return parsed, nil
}

// UUID extracts and validates a UUID path parameter.
//
// Validates RFC 4122 UUID format.
//
// Returns:
//   - The UUID value as string
//   - An error with detailed message if validation fails: "invalid {paramName} format: expected UUID, got '{value}'"
//
// Example:
//
//	userID, err := r.PathParam().UUID("id")
func (h PathParamHelper) UUID(key string) (string, error) {
	value := h.request.PathValue(key)
	if value == "" {
		return "", fmt.Errorf("missing path parameter: %s", key)
	}

	if !isValidUUID(value) {
		return "", fmt.Errorf("invalid %s format: expected UUID, got '%s'", key, value)
	}

	return value, nil
}

// isValidUUID validates if a string is a valid RFC 4122 UUID format.
// Format: 8-4-4-4-12 hex digits = 36 characters total.
func isValidUUID(s string) bool {
	// A UUID is 36 characters: 8-4-4-4-12 hex digits
	if len(s) != 36 {
		return false
	}

	// Check format with hyphens at correct positions
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return false
	}

	// Validate hex format for each segment
	segments := []string{s[0:8], s[9:13], s[14:18], s[19:23], s[24:36]}
	for _, segment := range segments {
		for _, ch := range segment {
			if !isHexChar(ch) {
				return false
			}
		}
	}

	return true
}

// isHexChar returns true if the character is a valid hex digit.
func isHexChar(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// PathParam returns a helper for extracting typed path parameters.
//
// Use it to extract typed path parameters:
//
//	userID, err := r.PathParam().String("id")
//	itemID, err := r.PathParam().Int64("itemId")
//	uuid, err := r.PathParam().UUID("id")
//
// Example:
//
//	userID, err := r.PathParam().String("id")
//	if err != nil {
//	    return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
//	}
func (r *Request[T]) PathParam() PathParamHelper {
	return r.pathParamHelper()
}
