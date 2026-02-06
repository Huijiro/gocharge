package gocharge_test

import (
	"context"
	"testing"
	"time"

	gocharge "github.com/huijiro/go-charge"
)

func TestWithValue(t *testing.T) {
	ctx := context.Background()

	// Test with string value
	ctx = gocharge.WithValue(ctx, "key", "value")
	val, ok := gocharge.Value[string](ctx, "key")

	if !ok {
		t.Error("Expected value to be found in context")
	}

	if val != "value" {
		t.Errorf("Expected 'value', got %s", val)
	}
}

func TestWithValueNotFound(t *testing.T) {
	ctx := context.Background()

	val, ok := gocharge.Value[string](ctx, "nonexistent")

	if ok {
		t.Error("Expected value to not be found in context")
	}

	if val != "" {
		t.Errorf("Expected empty string, got %s", val)
	}
}

func TestWithUser(t *testing.T) {
	ctx := context.Background()

	user := map[string]interface{}{"id": "123", "name": "Test User"}
	ctx = gocharge.WithUser(ctx, user)

	retrievedUser, ok := gocharge.UserFromContext(ctx)
	if !ok {
		t.Error("Expected user to be found in context")
	}

	userMap, ok := retrievedUser.(map[string]interface{})
	if !ok {
		t.Error("Expected user to be a map")
	}

	if userMap["id"] != "123" {
		t.Errorf("Expected id '123', got %v", userMap["id"])
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()

	traceID := "trace-123"
	ctx = gocharge.WithTraceID(ctx, traceID)

	retrievedTraceID, ok := gocharge.TraceIDFromContext(ctx)
	if !ok {
		t.Error("Expected trace ID to be found in context")
	}

	if retrievedTraceID != traceID {
		t.Errorf("Expected trace ID %s, got %s", traceID, retrievedTraceID)
	}
}

func TestContextTimeout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Check context is valid
	select {
	case <-ctx.Done():
		t.Error("Context should not be done yet")
	default:
		// Good
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	select {
	case <-ctx.Done():
		// Good
	default:
		t.Error("Context should have timed out")
	}
}
