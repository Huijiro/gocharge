package gocharge

import (
	"context"
	"encoding/json"
	"net/http"
)

// Request[T] wraps an http.Request with typed body data
type Request[T any] struct {
	Data T
	http.Request
}

// JSON decodes the request body as JSON into the Data field
func (r *Request[T]) JSON() (*T, error) {
	err := json.NewDecoder(r.Body).Decode(&r.Data)

	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}

// Context returns the request's context
// This is a convenience method to access the context from the embedded http.Request
func (r *Request[T]) Context() context.Context {
	return r.Request.Context()
}
