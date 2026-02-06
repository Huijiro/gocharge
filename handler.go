package gocharge

import (
	"context"
	"encoding/json"
	"net/http"
)

// HandlerFunc[W, R] is a typed HTTP request handler for gocharge.
//
// Type parameters:
//   - W: The response body type (JSON-encodable)
//   - R: The request body type (JSON-decodable)
//
// The handler receives:
//   - ctx: Request context for scoped values and deadlines
//   - w: Response wrapper for sending typed responses
//   - r: Request wrapper for accessing typed request data
//
// Returns an error which is automatically encoded by the server's ErrorEncoder.
// If the handler returns nil, the response has already been written.
// If it returns an error, the response will be replaced with an error response
// and the handler's response will be discarded.
//
// Example:
//
//	type CreateUserRequest struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	type UserResponse struct {
//	    ID    string `json:"id"`
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	func createUserHandler(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[CreateUserRequest]) error {
//	    req, err := r.JSON()
//	    if err != nil {
//	        return gocharge.NewAppError(gocharge.ErrBadRequest, "Invalid JSON")
//	    }
//	    if req.Name == "" {
//	        return gocharge.NewAppError(gocharge.ErrValidation, "Name is required")
//	    }
//	    user := UserResponse{ID: "new-id", Name: req.Name, Email: req.Email}
//	    _, err = w.Status(gocharge.StatusCreated).JSON(user)
//	    return err
//	}
type HandlerFunc[W any, R any] func(ctx context.Context, w Response[W], r Request[R]) error

// RegisterHandler registers a typed HTTP handler with the server.
//
// The handler is registered at the given path and wrapped to:
//   - Extract context from the HTTP request
//   - Extract path parameters and populate Request.PathParams
//   - Create typed request/response wrappers
//   - Call the typed handler function
//   - Catch and encode any returned errors
//   - Log handler registration
//
// Path can be any pattern supported by http.ServeMux:
//   - "/api/users" - simple path
//   - "/api/users/{id}" - with path segment
//   - "/api/v1/organizations/{id}/projects/{pid}/tasks" - multiple segments
//   - "GET /api/users/{id}" - with HTTP method (Go 1.22+)
//
// Path parameters are automatically extracted and accessible via:
//   - r.PathParam().String("id") - extract string parameter
//   - r.PathParam().Int("id") - extract integer parameter
//   - r.PathParam().Int64("id") - extract int64 parameter
//   - r.PathParam().UUID("id") - extract and validate UUID parameter
//
// Errors returned from the handler are automatically encoded by the server's
// ErrorEncoder and sent as JSON responses with appropriate HTTP status codes.
//
// Example:
//
//	gocharge.RegisterHandler(server, "/api/users", userHandler)
//	gocharge.RegisterHandler(server, "/api/users/{id}", getUserHandler)
//	gocharge.RegisterHandler(server, "GET /api/users/{id}", getDetailedUserHandler)
func RegisterHandler[W any, R any](s *Server, path string, handler HandlerFunc[W, R]) {
	// Log registration using server's logger
	s.Logger.Info(context.Background(), "registering_handler", "path", path)

	// Create a basic http.HandlerFunc that wraps our typed handler
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		// Get context from request
		ctx := r.Context()

		// Create request wrapper
		// Path parameters are accessed via r.PathParam().String/Int/Int64/UUID()
		// which uses Go 1.22+ http.Request.PathValue() internally
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
