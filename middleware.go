package gocharge

import (
	"context"
	"net/http"
)

// Middleware is a function that wraps an HTTP handler
// It receives a handler and returns a new handler with additional behavior
// Middleware executes in the order they are registered (first registered = outermost)
type Middleware func(next http.Handler) http.Handler

// Chain manages a chain of middleware
// Middleware is applied in registration order, so the first registered
// middleware is the outermost and executes first on request
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new middleware chain
func NewChain() *Chain {
	return &Chain{
		middlewares: make([]Middleware, 0),
	}
}

// Use adds one or more middleware to the chain
// Middleware added later will be closer to the actual handler
func (c *Chain) Use(middlewares ...Middleware) *Chain {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

// Build wraps the given handler with all registered middleware
// Returns a handler ready to be registered with http.ServeMux
func (c *Chain) Build(handler http.Handler) http.Handler {
	// Apply middleware in reverse order so first registered is outermost
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}

// BuildFunc is like Build but for HandlerFunc types
func (c *Chain) BuildFunc(handler http.HandlerFunc) http.Handler {
	return c.Build(handler)
}

// ContextMiddleware is a helper to create middleware that works with context
// This is useful for middleware that needs to pass data through context
func ContextMiddleware(fn func(ctx context.Context, next http.Handler, w http.ResponseWriter, r *http.Request)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			fn(ctx, next, w, r)
		})
	}
}
