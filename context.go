package gocharge

import (
	"context"
	"time"
)

// contextKey is used for storing values in context to avoid key collisions
type contextKey string

// Standard context keys used by the framework
const (
	contextKeyUser    contextKey = "gocharge:user"
	contextKeyTraceID contextKey = "gocharge:trace_id"
)

// ContextFromRequest creates a new context from an http.Request with optional timeout
// If timeout is 0, uses the request's context as-is
func ContextFromRequest(r *Request[interface{}], timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := r.Context()

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}

	// Return a no-op cancel function if no timeout
	return ctx, func() {}
}

// WithValue sets a generic value in context
// This is a type-safe wrapper around context.WithValue
func WithValue[T any](ctx context.Context, key string, value T) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

// Value retrieves a typed value from context
// Returns the value and a boolean indicating if it was found
func Value[T any](ctx context.Context, key string) (T, bool) {
	val := ctx.Value(contextKey(key))
	if val == nil {
		var zero T
		return zero, false
	}

	if typedVal, ok := val.(T); ok {
		return typedVal, true
	}

	var zero T
	return zero, false
}

// WithUser sets the user in context
func WithUser(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, contextKeyUser, user)
}

// UserFromContext retrieves the user from context
func UserFromContext(ctx context.Context) (interface{}, bool) {
	return Value[interface{}](ctx, string(contextKeyUser))
}

// WithTraceID sets the trace ID in context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, contextKeyTraceID, traceID)
}

// TraceIDFromContext retrieves the trace ID from context
func TraceIDFromContext(ctx context.Context) (string, bool) {
	return Value[string](ctx, string(contextKeyTraceID))
}
