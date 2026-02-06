package gocharge_test

import (
	"context"
	"net/http"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

func TestErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		code     gocharge.ErrorCode
		expected string
	}{
		{
			name:     "BadRequest",
			code:     gocharge.ErrBadRequest,
			expected: "BAD_REQUEST",
		},
		{
			name:     "NotFound",
			code:     gocharge.ErrNotFound,
			expected: "NOT_FOUND",
		},
		{
			name:     "InternalError",
			code:     gocharge.ErrInternalServerError,
			expected: "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.code) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.code))
			}
		})
	}
}

func TestAppError(t *testing.T) {
	err := gocharge.NewError(gocharge.ErrNotFound, "Resource not found", http.StatusNotFound)

	if err.Error() != "Resource not found" {
		t.Errorf("Expected 'Resource not found', got %s", err.Error())
	}

	if err.Code() != gocharge.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err.Code())
	}

	if err.StatusCode() != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", err.StatusCode())
	}
}

func TestNewAppError(t *testing.T) {
	tests := []struct {
		name           string
		code           gocharge.ErrorCode
		expectedStatus int
	}{
		{
			name:           "BadRequest",
			code:           gocharge.ErrBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NotFound",
			code:           gocharge.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unauthorized",
			code:           gocharge.ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "InternalServerError",
			code:           gocharge.ErrInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gocharge.NewAppError(tt.code, "test error")
			if err.StatusCode() != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, err.StatusCode())
			}
		})
	}
}

func TestDefaultErrorEncoder(t *testing.T) {
	encoder := &gocharge.DefaultErrorEncoder{}
	ctx := context.Background()

	// Test with Errorable error
	appErr := gocharge.NewAppError(gocharge.ErrNotFound, "Not found")
	errResp, status := encoder.Encode(ctx, appErr)

	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}

	if errResp.Code != string(gocharge.ErrNotFound) {
		t.Errorf("Expected code NOT_FOUND, got %s", errResp.Code)
	}

	if errResp.Message != "Not found" {
		t.Errorf("Expected message 'Not found', got %s", errResp.Message)
	}

	// Test with generic error
	genericErr := errorStub("generic error")
	errResp, status = encoder.Encode(ctx, genericErr)

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", status)
	}

	if errResp.Code != string(gocharge.ErrInternalServerError) {
		t.Errorf("Expected code INTERNAL_SERVER_ERROR, got %s", errResp.Code)
	}
}

// errorStub is a simple error implementation for testing
type errorStub string

func (e errorStub) Error() string {
	return string(e)
}
