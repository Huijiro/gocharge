package gocharge

import (
	"context"
	"encoding/json"
	"net/http"
)

// HandlerFunc[W, R] is the handler function signature for gocharge
// W is the response body type, R is the request body type
// Handlers receive context as the first parameter and must return an error or nil
type HandlerFunc[W any, R any] func(ctx context.Context, w Response[W], r Request[R]) error

// RegisterHandler registers a typed handler with the server
// The handler will be wrapped to handle errors, context, and middleware
func RegisterHandler[W any, R any](s *Server, path string, handler HandlerFunc[W, R]) {
	// Log registration using server's logger
	s.Logger.Info(context.Background(), "registering_handler", "path", path)

	// Create a basic http.HandlerFunc that wraps our typed handler
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		// Get context from request
		ctx := r.Context()

		// Create request wrapper
		request := Request[R]{
			Request: *r,
			Data:    *new(R),
		}

		// Create response wrapper
		response := Response[W]{
			ResponseWriter: w,
			Data:           *new(W),
			StatusCode:     0, // Will default to 200 if not set
		}

		// Call the handler and capture any error
		err := handler(ctx, response, request)

		// If there was an error, encode it using the server's error encoder
		if err != nil {
			errResp, statusCode := s.ErrorEncoder.Encode(ctx, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)

			// Encode the error response
			if encodeErr := json.NewEncoder(w).Encode(errResp); encodeErr != nil {
				s.Logger.Error(ctx, "failed_to_encode_error_response", encodeErr)
			}
		}
	}

	// Register the handler with the server's mux
	s.Handler.HandleFunc(path, httpHandler)
}
