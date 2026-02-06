// Package gocharge provides a type-safe, minimal HTTP server framework for building
// production-ready Go backends with zero external dependencies.
//
// # Overview
//
// GoCharge is designed around these core principles:
//   - Type safety: Leverage Go 1.18+ generics for request/response types
//   - Error handling: Proper error propagation with typed error codes
//   - Context propagation: Built-in context support for request scoping
//   - Structured logging: Using Go 1.21+ slog for production logging
//   - Middleware composition: Composable middleware chain system
//   - Stdlib-only: Uses only Go standard library, no external dependencies
//
// # Quick Start
//
// Create a server and register typed handlers:
//
//	type UserResponse struct {
//	    ID   string `json:"id"`
//	    Name string `json:"name"`
//	}
//
//	func helloHandler(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[string]) error {
//	    _, err := w.JSON(UserResponse{ID: "1", Name: "Alice"})
//	    return err
//	}
//
//	server := gocharge.New(":8080")
//	gocharge.RegisterHandler(server, "/api/users", helloHandler)
//	server.ListenAndServe()
//
// # Error Handling
//
// Handlers return typed errors that are automatically encoded to JSON:
//
//	func getUserHandler(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[GetRequest]) error {
//	    if req.ID == "" {
//	        return gocharge.NewAppError(gocharge.ErrValidation, "ID is required")
//	    }
//	    user, found := database.GetUser(req.ID)
//	    if !found {
//	        return gocharge.NewAppError(gocharge.ErrNotFound, "User not found")
//	    }
//	    _, err := w.JSON(user)
//	    return err
//	}
//
// # Response Chaining
//
// Response methods support chaining for fluent API:
//
//	_, err := w.Status(gocharge.StatusCreated).JSON(newUser)
//	return err
//
// # Context and Request Scoping
//
// Context flows through the entire request lifecycle and can be used to store
// request-scoped values like user info or trace IDs:
//
//	func protectedHandler(ctx context.Context, w gocharge.Response[Response], r gocharge.Request[Request]) error {
//	    ctx = gocharge.WithUser(ctx, user)
//	    logger := gocharge.GetLogger()
//	    logger.Info(ctx, "processing request")
//	    _, err := w.JSON(response)
//	    return err
//	}
//
// # Structured Logging
//
// GoCharge uses Go 1.21+ slog for structured logging:
//
//	logger := gocharge.GetLogger()
//	logger.Info(ctx, "user created", "user_id", userID)
//	logger.Error(ctx, "failed to save user", err)
//
// Users can provide custom loggers by implementing the Logger interface.
//
// # Middleware
//
// Middleware can be registered to apply cross-cutting concerns:
//
//	chain := server.Middleware()
//	chain.Use(middleware.LoggingMiddleware)
//	chain.Use(middleware.RecoveryMiddleware)
//
// See package middleware for built-in middleware implementations.
//
// # Sub-routes and API Organization
//
// Any path pattern is supported through http.ServeMux:
//
//	// Basic routes
//	gocharge.RegisterHandler(server, "/api/users", listUsers)
//	gocharge.RegisterHandler(server, "/api/users/{id}", getUser)
//
//	// Versioned routes
//	gocharge.RegisterHandler(server, "/api/v1/users", v1ListUsers)
//	gocharge.RegisterHandler(server, "/api/v2/users", v2ListUsers)
//
//	// RESTful nested resources
//	gocharge.RegisterHandler(server, "/api/users/{id}/posts", getUserPosts)
//
// See SUBROUTES.md for comprehensive patterns.
package gocharge

import (
	"net/http"
	"time"
)

// Server is the main HTTP server for gocharge applications.
//
// It extends http.Server with typed error handling, structured logging,
// and middleware chain management. All handlers registered with this server
// receive context, typed request/response wrappers, and automatic error
// encoding.
//
// Fields:
//   - Handler: The underlying http.ServeMux for route registration
//   - ErrorEncoder: Encodes errors to ErrorResponse JSON (customizable)
//   - Logger: Structured logger for application logging (customizable)
//   - MiddlewareChain: Composable middleware chain (optional)
type Server struct {
	http.Server
	Handler         *http.ServeMux
	ErrorEncoder    ErrorEncoder
	Logger          Logger
	MiddlewareChain *Chain
}

// New creates a new Server listening on the given address.
//
// The server is initialized with:
//   - DefaultErrorEncoder for standard error responses
//   - Default structured logger using slog
//   - Empty middleware chain (optional)
//
// Example:
//
//	server := gocharge.New(":8080")
//	gocharge.RegisterHandler(server, "/api/users", handler)
//	server.ListenAndServe()
func New(addr string) *Server {
	mux := http.NewServeMux()
	return &Server{
		Server: http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 15 * time.Second,
		},
		Handler:         mux,
		ErrorEncoder:    &DefaultErrorEncoder{},
		Logger:          GetLogger(),
		MiddlewareChain: NewChain(),
	}
}

// SetErrorEncoder sets a custom error encoder for the server.
//
// The error encoder is called whenever a handler returns an error.
// It converts the error to an ErrorResponse and determines the HTTP status code.
// This allows customizing how errors are presented to clients.
//
// Example:
//
//	type CustomErrorEncoder struct{}
//	func (e *CustomErrorEncoder) Encode(ctx context.Context, err error) (ErrorResponse, int) {
//	    // Custom error encoding logic
//	}
//	server.SetErrorEncoder(&CustomErrorEncoder{})
func (s *Server) SetErrorEncoder(enc ErrorEncoder) {
	s.ErrorEncoder = enc
}

// SetLogger sets the structured logger for the server.
//
// The logger is used throughout the request lifecycle for logging events.
// Users can implement the Logger interface to provide custom logging behavior.
// By default, uses a slog logger outputting JSON to stdout.
//
// Example:
//
//	type CustomLogger struct{}
//	func (l *CustomLogger) Info(ctx context.Context, msg string, args ...any) { }
//	server.SetLogger(&CustomLogger{})
func (s *Server) SetLogger(l Logger) {
	s.Logger = l
}

// Middleware returns the middleware chain for this server.
//
// The chain can be used to register middleware that applies to all handlers.
// Middleware is executed in registration order (first registered = outermost).
//
// Example:
//
//	chain := server.Middleware()
//	chain.Use(middleware.LoggingMiddleware)
//	chain.Use(middleware.RecoveryMiddleware)
func (s *Server) Middleware() *Chain {
	return s.MiddlewareChain
}
