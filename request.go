package gocharge

import (
	"context"
	"encoding/json"
	"net/http"
)

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
