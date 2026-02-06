package gocharge

import (
	"context"
	"log/slog"
	"os"
)

// Logger is the interface for structured logging in gocharge
// Users can implement this to provide custom logging behavior
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, err error, args ...any)
}

// DefaultLogger is the default implementation using log/slog
type DefaultLogger struct {
	logger *slog.Logger
}

// NewDefaultLogger creates a new default logger using slog
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// NewDefaultLoggerWithHandler creates a default logger with a custom slog handler
func NewDefaultLoggerWithHandler(handler slog.Handler) *DefaultLogger {
	return &DefaultLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a debug-level message
func (l *DefaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// Info logs an info-level message
func (l *DefaultLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// Warn logs a warn-level message
func (l *DefaultLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// Error logs an error-level message with the error
func (l *DefaultLogger) Error(ctx context.Context, msg string, err error, args ...any) {
	l.logger.ErrorContext(ctx, msg, append([]any{"error", err}, args...)...)
}

// Global logger instance
var globalLogger Logger = NewDefaultLogger()

// SetLogger sets the global logger instance
func SetLogger(l Logger) {
	globalLogger = l
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	return globalLogger
}
