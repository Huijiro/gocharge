package gocharge_test

import (
	"context"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

// MockLogger is a test implementation of Logger
type MockLogger struct {
	debugCalls int
	infoCalls  int
	warnCalls  int
	errorCalls int
}

func (m *MockLogger) Debug(ctx context.Context, msg string, args ...any) {
	m.debugCalls++
}

func (m *MockLogger) Info(ctx context.Context, msg string, args ...any) {
	m.infoCalls++
}

func (m *MockLogger) Warn(ctx context.Context, msg string, args ...any) {
	m.warnCalls++
}

func (m *MockLogger) Error(ctx context.Context, msg string, err error, args ...any) {
	m.errorCalls++
}

func TestDefaultLogger(t *testing.T) {
	logger := gocharge.NewDefaultLogger()
	ctx := context.Background()

	// Just verify these don't panic
	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, "error message", errorStub("test error"))
}

func TestSetLogger(t *testing.T) {
	original := gocharge.GetLogger()
	defer gocharge.SetLogger(original)

	mockLogger := &MockLogger{}
	gocharge.SetLogger(mockLogger)

	logger := gocharge.GetLogger()
	ctx := context.Background()

	logger.Info(ctx, "test message")
	if mockLogger.infoCalls != 1 {
		t.Errorf("Expected 1 info call, got %d", mockLogger.infoCalls)
	}

	logger.Debug(ctx, "debug message")
	if mockLogger.debugCalls != 1 {
		t.Errorf("Expected 1 debug call, got %d", mockLogger.debugCalls)
	}

	logger.Error(ctx, "error", errorStub("test"))
	if mockLogger.errorCalls != 1 {
		t.Errorf("Expected 1 error call, got %d", mockLogger.errorCalls)
	}
}

func TestGetLogger(t *testing.T) {
	logger := gocharge.GetLogger()
	if logger == nil {
		t.Error("Expected logger to be non-nil")
	}
}
