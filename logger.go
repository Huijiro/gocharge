package gocharge

import (
	"context"
	"log/slog"
	"os"
)

// Logger defines the structured logging interface for gocharge.
//
// Applications can implement this interface to provide custom logging behavior.
// The logger is passed context for request-scoped values like trace IDs.
//
// The args parameter follows slog conventions: alternating keys and values.
//
// Example custom implementation:
//
//	type MyLogger struct{}
//	func (l *MyLogger) Info(ctx context.Context, msg string, args ...any) {
//	    // Custom logging logic
//	}
type Logger interface {
	// Debug logs a debug-level message with optional key-value pairs
	Debug(ctx context.Context, msg string, args ...any)
	// Info logs an info-level message with optional key-value pairs
	Info(ctx context.Context, msg string, args ...any)
	// Warn logs a warning-level message with optional key-value pairs
	Warn(ctx context.Context, msg string, args ...any)
	// Error logs an error-level message with the error and optional key-value pairs
	Error(ctx context.Context, msg string, err error, args ...any)
}

// DefaultLogger is a built-in Logger implementation using Go's log/slog.
//
// It outputs JSON-formatted structured logs to stdout by default.
// The logger respects the slog global log level configuration.
type DefaultLogger struct {
	logger *slog.Logger
}

// NewDefaultLogger creates a new logger using slog with JSON output.
//
// Logs are sent to os.Stdout in JSON format, compatible with log aggregation
// services like ELK, Splunk, and cloud providers' log systems.
//
// Example:
//
//	logger := gocharge.NewDefaultLogger()
//	gocharge.SetLogger(logger)
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// NewDefaultLoggerWithHandler creates a logger with a custom slog handler.
//
// This allows customizing the log format, output destination, or level.
//
// Example with text output:
//
//	handler := slog.NewTextHandler(os.Stdout, nil)
//	logger := gocharge.NewDefaultLoggerWithHandler(handler)
//
// Example with custom options:
//
//	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
//	handler := slog.NewJSONHandler(os.Stderr, opts)
//	logger := gocharge.NewDefaultLoggerWithHandler(handler)
func NewDefaultLoggerWithHandler(handler slog.Handler) *DefaultLogger {
	return &DefaultLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a debug-level message.
//
// Debug logs are typically disabled in production.
// Args follow slog conventions: alternating keys and values.
func (l *DefaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// Info logs an info-level message.
//
// Info logs are typically enabled in production for important events.
// Args follow slog conventions: alternating keys and values.
func (l *DefaultLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// Warn logs a warning-level message.
//
// Warning logs indicate potentially problematic but non-critical events.
// Args follow slog conventions: alternating keys and values.
func (l *DefaultLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// Error logs an error-level message with an associated error.
//
// Error logs should be used for recoverable errors. The error is automatically
// included in the log output.
// Additional args follow slog conventions: alternating keys and values.
func (l *DefaultLogger) Error(ctx context.Context, msg string, err error, args ...any) {
	l.logger.ErrorContext(ctx, msg, append([]any{"error", err}, args...)...)
}

// globalLogger is the package-level logger used by the framework.
var globalLogger Logger = NewDefaultLogger()

// SetLogger sets the global logger used by the framework.
//
// Call this once at application startup to use a custom logger.
// All framework and application code should use GetLogger() to access the logger.
//
// Example:
//
//	gocharge.SetLogger(&MyCustomLogger{})
func SetLogger(l Logger) {
	globalLogger = l
}

// GetLogger returns the global logger for this application.
//
// This should be called in handler code and middleware to log events.
//
// Example:
//
//	func handler(ctx context.Context, w Response[R], r Request[Req]) error {
//	    logger := gocharge.GetLogger()
//	    logger.Info(ctx, "processing request", "id", id)
//	    return nil
//	}
func GetLogger() Logger {
	return globalLogger
}
