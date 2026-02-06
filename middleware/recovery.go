package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/huijiro/go-charge"
)

// RecoveryMiddleware recovers from panics and returns a 500 error response
// This prevents the server from crashing due to unhandled panics in handlers
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := r.Context()
				logger := gocharge.GetLogger()

				// Log the panic with stack trace
				logger.Error(ctx, "panic_recovered",
					fmt.Errorf("panic: %v", err),
					"panic_value", fmt.Sprintf("%v", err),
					"stack_trace", string(debug.Stack()),
				)

				// Write error response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				errResp := gocharge.ErrorResponse{
					Code:    string(gocharge.ErrInternalServerError),
					Message: "An internal server error occurred",
				}

				// Try to encode error response, but don't panic if encoding fails
				if err := encodeErrorResponse(w, errResp); err != nil {
					logger.Error(ctx, "failed_to_encode_error_response", err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// encodeErrorResponse is a helper to encode error responses
func encodeErrorResponse(w http.ResponseWriter, errResp gocharge.ErrorResponse) error {
	return json.NewEncoder(w).Encode(errResp)
}
