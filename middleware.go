package gocharge

import (
	"context"
	"net/http"
)

// Middleware is a function that wraps an HTTP handler with additional behavior.
//
// Middleware receives a handler and returns a new handler. This allows for
// composable cross-cutting concerns like logging, authentication, and error recovery.
//
// Middleware executes in registration order: the first registered middleware
// is the outermost (executes first on request, last on response).
//
// Example middleware that adds context values:
//
//	func AuthMiddleware(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        token := r.Header.Get("Authorization")
//	        if token == "" {
//	            http.Error(w, "Unauthorized", http.StatusUnauthorized)
//	            return
//	        }
//	        ctx := gocharge.WithUser(r.Context(), token)
//	        next.ServeHTTP(w, r.WithContext(ctx))
//	    })
//	}
type Middleware func(next http.Handler) http.Handler

// Chain manages an ordered collection of middleware.
//
// Middleware is applied in registration order, so the first middleware registered
// will wrap the handler first (outermost) and execute first on request.
//
// Example:
//
//	chain := gocharge.NewChain()
//	chain.Use(LoggingMiddleware)      // Outermost - logs all requests
//	chain.Use(RecoveryMiddleware)     // Catches panics
//	chain.Use(AuthMiddleware)         // Authenticates requests
//	                                  // Inner: actual handler
//
// Request flow: Request → Logging → Recovery → Auth → Handler → Auth → Recovery → Logging → Response
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new, empty middleware chain.
//
// Example:
//
//	chain := gocharge.NewChain()
//	chain.Use(myMiddleware)
func NewChain() *Chain {
	return &Chain{
		middlewares: make([]Middleware, 0),
	}
}

// Use adds one or more middleware to the chain.
//
// Middleware is added in the order specified. Each middleware wraps the ones added before it,
// so the first call to Use() creates the outermost middleware.
//
// Returns the receiver for method chaining.
//
// Example:
//
//	chain.Use(middleware.LoggingMiddleware, middleware.RecoveryMiddleware)
//	// or
//	chain.Use(middleware.LoggingMiddleware).Use(middleware.RecoveryMiddleware)
func (c *Chain) Use(middlewares ...Middleware) *Chain {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

// Build wraps an HTTP handler with all registered middleware.
//
// Middleware is applied in reverse registration order so the first registered
// middleware becomes the outermost layer (first to execute on request).
//
// Returns an http.Handler ready to be used with http.HandleFunc or http.Handle.
//
// Once a handler is built, it can be used directly:
//
//	handler := chain.Build(myHandler)
//	server.Handler.Handle("/path", handler)
//
// The registered handlers use this automatically, so in normal usage you don't need to call Build.
func (c *Chain) Build(handler http.Handler) http.Handler {
	// Apply middleware in reverse order so first registered is outermost
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}

// BuildFunc is a convenience method for wrapping http.HandlerFunc.
//
// Equivalent to Build but accepts http.HandlerFunc instead of http.Handler.
func (c *Chain) BuildFunc(handler http.HandlerFunc) http.Handler {
	return c.Build(handler)
}

// ContextMiddleware is a helper for creating middleware that works with context.
//
// This is a convenience for middleware that needs to work with context values.
// The provided function receives the context, next handler, and response/request.
//
// Example:
//
//	middleware := gocharge.ContextMiddleware(func(ctx context.Context, next http.Handler, w http.ResponseWriter, r *http.Request) {
//	    traceID := "trace-" + generateID()
//	    ctx = gocharge.WithTraceID(ctx, traceID)
//	    next.ServeHTTP(w, r.WithContext(ctx))
//	})
func ContextMiddleware(fn func(ctx context.Context, next http.Handler, w http.ResponseWriter, r *http.Request)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			fn(ctx, next, w, r)
		})
	}
}
