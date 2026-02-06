package middleware

import (
	"net/http"
	"time"

	"github.com/huijiro/go-charge"
)

// LoggingMiddleware logs HTTP requests and responses
// It logs the request method, path, and response status code and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := gocharge.GetLogger()

		// Log request start
		logger.Info(ctx, "request_start",
			"method", r.Method,
			"path", r.RequestURI,
			"remote_addr", r.RemoteAddr,
		)

		// Create a response writer that captures status code
		wrapped := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Record start time
		startTime := time.Now()

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log request end with duration and status
		duration := time.Since(startTime)
		logger.Info(ctx, "request_end",
			"method", r.Method,
			"path", r.RequestURI,
			"status_code", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures the status code
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write implements http.ResponseWriter.Write
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}
