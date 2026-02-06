package gocharge_test

import (
	"net/http"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

func TestChainCreation(t *testing.T) {
	chain := gocharge.NewChain()
	if chain == nil {
		t.Error("Expected chain to be non-nil")
	}
}

func TestChainUse(t *testing.T) {
	chain := gocharge.NewChain()

	// Create dummy middleware
	middleware1 := func(next http.Handler) http.Handler {
		return next
	}

	middleware2 := func(next http.Handler) http.Handler {
		return next
	}

	// Test adding middleware
	chain.Use(middleware1)
	chain.Use(middleware2)

	// Verify chain doesn't panic on build
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	builtHandler := chain.Build(handler)
	if builtHandler == nil {
		t.Error("Expected built handler to be non-nil")
	}
}

func TestChainOrder(t *testing.T) {
	// This test verifies that middleware executes in the correct order
	// First registered should be outermost (execute first)

	var executionOrder []string

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		})
	}

	chain := gocharge.NewChain()
	chain.Use(middleware1)
	chain.Use(middleware2)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	builtHandler := chain.Build(handler)

	// Create a mock response writer and request
	req, _ := http.NewRequest("GET", "/test", nil)
	w := &mockResponseWriter{header: make(http.Header)}

	builtHandler.ServeHTTP(w, req)

	// Expected order: middleware1-before, middleware2-before, handler, middleware2-after, middleware1-after
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Expected %d calls, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(executionOrder) {
			t.Errorf("Expected call %d to be %s, but no more calls", i, expected)
			break
		}
		if executionOrder[i] != expected {
			t.Errorf("Expected call %d to be %s, got %s", i, expected, executionOrder[i])
		}
	}
}

// mockResponseWriter is a simple mock for testing
type mockResponseWriter struct {
	header     http.Header
	statusCode int
	written    bool
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	if !m.written {
		m.statusCode = http.StatusOK
		m.written = true
	}
	return len(b), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	if !m.written {
		m.statusCode = statusCode
		m.written = true
	}
}
