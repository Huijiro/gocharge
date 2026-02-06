package gocharge

import (
	"context"
	"time"
)

// contextKey is a private type for context keys to avoid collisions.
//
// Using a typed key prevents other packages from accidentally (or intentionally)
// overwriting framework context values.
type contextKey string

// Standard context keys reserved by the framework for built-in values.
const (
	contextKeyUser    contextKey = "gocharge:user"
	contextKeyTraceID contextKey = "gocharge:trace_id"
)

// ContextFromRequest creates a context from an HTTP request with optional timeout.
//
// If timeout is 0, returns the request's context as-is with a no-op cancel function.
// If timeout > 0, wraps the context with context.WithTimeout.
//
// Returns a context and cancel function following the standard Go pattern.
// The cancel function should be deferred to ensure cleanup.
//
// Example:
//
//	ctx, cancel := gocharge.ContextFromRequest(r, 30*time.Second)
//	defer cancel()
//	result := someOperation(ctx)
func ContextFromRequest(r *Request[interface{}], timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := r.Context()

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}

	// Return a no-op cancel function if no timeout
	return ctx, func() {}
}

// WithValue stores a typed value in context with the given key.
//
// This is a type-safe wrapper around context.WithValue that preserves the type
// of the value. Use Value[T] to retrieve the value.
//
// Example:
//
//	ctx = gocharge.WithValue(ctx, "request_id", "req-123")
//	id, ok := gocharge.Value[string](ctx, "request_id")
func WithValue[T any](ctx context.Context, key string, value T) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

// Value retrieves a typed value from context.
//
// Returns the value and a boolean indicating if it was found.
// The zero value of T is returned if the key is not found or has a different type.
//
// Example:
//
//	requestID, ok := gocharge.Value[string](ctx, "request_id")
//	if !ok {
//	    requestID = "unknown"
//	}
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

// WithUser stores a user value in context.
//
// The user value can be any type (string, struct, interface{}, etc.).
// Use UserFromContext to retrieve the value.
//
// Example:
//
//	type User struct {
//	    ID   string
//	    Name string
//	}
//	ctx = gocharge.WithUser(ctx, User{ID: "123", Name: "Alice"})
//	user, ok := gocharge.UserFromContext(ctx)
func WithUser(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, contextKeyUser, user)
}

// UserFromContext retrieves the user from context.
//
// Returns the user value and a boolean indicating if it was found.
// The user value is returned as interface{} and should be type-asserted.
//
// Example:
//
//	user, ok := gocharge.UserFromContext(ctx)
//	if !ok {
//	    return gocharge.NewAppError(gocharge.ErrUnauthorized, "User not found")
//	}
//	userStruct := user.(User)
func UserFromContext(ctx context.Context) (interface{}, bool) {
	return Value[interface{}](ctx, string(contextKeyUser))
}

// WithTraceID stores a trace ID in context for request tracing.
//
// Trace IDs are useful for correlating logs and requests across services.
// Use TraceIDFromContext to retrieve the value.
//
// Example:
//
//	ctx = gocharge.WithTraceID(ctx, "trace-abc123")
//	traceID, ok := gocharge.TraceIDFromContext(ctx)
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, contextKeyTraceID, traceID)
}

// TraceIDFromContext retrieves the trace ID from context.
//
// Returns the trace ID and a boolean indicating if it was found.
//
// Example:
//
//	traceID, ok := gocharge.TraceIDFromContext(ctx)
//	if ok {
//	    logger.Info(ctx, "request", "trace_id", traceID)
//	}
func TraceIDFromContext(ctx context.Context) (string, bool) {
	return Value[string](ctx, string(contextKeyTraceID))
}
