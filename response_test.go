package gocharge_test

import (
	"net/http"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

type TestData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestResponseStatus(t *testing.T) {
	mockWriter := &mockResponseWriter{header: make(http.Header)}
	resp := gocharge.Response[TestData]{
		ResponseWriter: mockWriter,
		Data:           TestData{},
		StatusCode:     0,
	}

	// Test chaining
	result := resp.Status(gocharge.StatusCreated)

	if result != &resp {
		t.Error("Expected chaining to return the same response object")
	}

	if resp.StatusCode != int(gocharge.StatusCreated) {
		t.Errorf("Expected status code %d, got %d", int(gocharge.StatusCreated), resp.StatusCode)
	}
}

func TestResponseJSON(t *testing.T) {
	mockWriter := &mockResponseWriter{header: make(http.Header)}
	resp := gocharge.Response[TestData]{
		ResponseWriter: mockWriter,
		Data:           TestData{},
		StatusCode:     0,
	}

	testData := TestData{
		Message: "Test message",
		Code:    123,
	}

	// Test JSON method
	result, err := resp.JSON(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != &resp {
		t.Error("Expected chaining to return the same response object")
	}

	// StatusCode should remain 0 until WriteHeader is called
	// When StatusCode is 0, WriteHeader writes StatusOK
	if mockWriter.statusCode != http.StatusOK {
		t.Errorf("Expected written status code %d, got %d", http.StatusOK, mockWriter.statusCode)
	}
}

func TestResponseJSONWithStatus(t *testing.T) {
	mockWriter := &mockResponseWriter{header: make(http.Header)}
	resp := gocharge.Response[TestData]{
		ResponseWriter: mockWriter,
		Data:           TestData{},
		StatusCode:     0,
	}

	testData := TestData{
		Message: "Created",
		Code:    1,
	}

	// Chain status and JSON
	_, err := resp.Status(gocharge.StatusCreated).JSON(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockWriter.statusCode != int(gocharge.StatusCreated) {
		t.Errorf("Expected written status code %d, got %d", int(gocharge.StatusCreated), mockWriter.statusCode)
	}
}

func TestResponseError(t *testing.T) {
	mockWriter := &mockResponseWriter{header: make(http.Header)}
	resp := gocharge.Response[TestData]{
		ResponseWriter: mockWriter,
		Data:           TestData{},
		StatusCode:     0,
	}

	errorCode := "TEST_ERROR"
	errorMessage := "This is a test error"

	result, err := resp.Error(errorCode, errorMessage)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != &resp {
		t.Error("Expected chaining to return the same response object")
	}

	if mockWriter.statusCode != http.StatusInternalServerError {
		t.Errorf("Expected default error status code %d, got %d", http.StatusInternalServerError, mockWriter.statusCode)
	}

	// The content should have been written
	if mockWriter.header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected content type application/json, got %s", mockWriter.header.Get("Content-Type"))
	}
}

func TestResponseErrorWithStatus(t *testing.T) {
	mockWriter := &mockResponseWriter{header: make(http.Header)}
	resp := gocharge.Response[TestData]{
		ResponseWriter: mockWriter,
		Data:           TestData{},
		StatusCode:     0,
	}

	// Chain status and error
	_, err := resp.Status(gocharge.StatusBadRequest).Error("VALIDATION_ERROR", "Invalid input")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockWriter.statusCode != int(gocharge.StatusBadRequest) {
		t.Errorf("Expected written status code %d, got %d", int(gocharge.StatusBadRequest), mockWriter.statusCode)
	}
}

func TestResponseIntegration(t *testing.T) {
	// Test the full flow with actual JSON encoding
	mockWriter := &mockResponseWriter{header: make(http.Header)}
	resp := gocharge.Response[TestData]{
		ResponseWriter: mockWriter,
		Data:           TestData{},
		StatusCode:     0,
	}

	data := TestData{
		Message: "Success",
		Code:    200,
	}

	_, err := resp.Status(gocharge.StatusOK).JSON(data)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the data was set
	if resp.Data.Message != "Success" {
		t.Errorf("Expected message 'Success', got %s", resp.Data.Message)
	}
}
