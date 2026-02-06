package gocharge_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

// TestErrorHandling verifies that handler errors are properly caught and encoded
func TestErrorHandling(t *testing.T) {
	type ErrorTestResponse struct {
		Status string `json:"status"`
	}

	gocharge.RegisterHandler(server, "/testError", func(ctx context.Context, w gocharge.Response[ErrorTestResponse], r gocharge.Request[string]) error {
		// Return a typed error
		return gocharge.NewAppError(gocharge.ErrNotFound, "Resource not found")
	})

	response, err := http.Get("http://localhost:8080/testError")
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	// Verify error status code
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", response.StatusCode)
	}

	// Verify error response body
	var errResp gocharge.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&errResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if errResp.Code != string(gocharge.ErrNotFound) {
		t.Errorf("Expected error code NOT_FOUND, got %s", errResp.Code)
	}

	if errResp.Message != "Resource not found" {
		t.Errorf("Expected message 'Resource not found', got %s", errResp.Message)
	}
}

// TestContextPropagation verifies that context is properly passed to handlers
func TestContextPropagation(t *testing.T) {
	type ContextTestResponse struct {
		HasUser bool `json:"has_user"`
	}

	gocharge.RegisterHandler(server, "/testContext", func(ctx context.Context, w gocharge.Response[ContextTestResponse], r gocharge.Request[string]) error {
		// Try to get user from context (won't be there unless middleware sets it)
		_, hasUser := gocharge.UserFromContext(ctx)

		resp := ContextTestResponse{HasUser: hasUser}
		_, err := w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/testContext")
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	var ctxResp ContextTestResponse
	err = json.NewDecoder(response.Body).Decode(&ctxResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	// Should be false since no middleware set it
	if ctxResp.HasUser {
		t.Error("Expected HasUser to be false")
	}
}

// TestResponseChaining verifies that response methods can be chained
func TestResponseChaining(t *testing.T) {
	type ChainTestResponse struct {
		Message string `json:"message"`
	}

	gocharge.RegisterHandler(server, "/testChain", func(ctx context.Context, w gocharge.Response[ChainTestResponse], r gocharge.Request[string]) error {
		// Chain Status and JSON
		_, err := w.Status(gocharge.StatusCreated).JSON(ChainTestResponse{Message: "Created"})
		return err
	})

	response, err := http.Get("http://localhost:8080/testChain")
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	// Verify status code was set by chaining
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", response.StatusCode)
	}

	var chainResp ChainTestResponse
	err = json.NewDecoder(response.Body).Decode(&chainResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if chainResp.Message != "Created" {
		t.Errorf("Expected message 'Created', got %s", chainResp.Message)
	}
}

// TestValidationError demonstrates custom error handling
func TestValidationError(t *testing.T) {
	type ValidateRequest struct {
		Name string `json:"name"`
	}

	type ValidateResponse struct {
		Success bool `json:"success"`
	}

	gocharge.RegisterHandler(server, "/testValidate", func(ctx context.Context, w gocharge.Response[ValidateResponse], r gocharge.Request[ValidateRequest]) error {
		body, err := r.JSON()
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrValidation, "Invalid JSON body")
		}

		if body.Name == "" {
			return gocharge.NewAppError(gocharge.ErrValidation, "Name is required")
		}

		_, err = w.JSON(ValidateResponse{Success: true})
		return err
	})

	// Test with empty body
	response, err := http.Get("http://localhost:8080/testValidate")
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty request, got %d", response.StatusCode)
	}

	// Test with invalid JSON
	invalidBody := bytes.NewBufferString("not json")
	request, _ := http.NewRequest("POST", "http://localhost:8080/testValidate", invalidBody)
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", response.StatusCode)
	}

	// Test with valid data
	validBody := bytes.NewBufferString(`{"name": "Test User"}`)
	request, _ = http.NewRequest("POST", "http://localhost:8080/testValidate", validBody)
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for valid request, got %d", response.StatusCode)
	}

	var validateResp ValidateResponse
	err = json.NewDecoder(response.Body).Decode(&validateResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if !validateResp.Success {
		t.Error("Expected success to be true")
	}
}

// TestErrorResponse verifies the Error method on Response
func TestErrorResponse(t *testing.T) {
	type ErrorMethodResponse struct {
		Unused string `json:"unused"`
	}

	gocharge.RegisterHandler(server, "/testErrorMethod", func(ctx context.Context, w gocharge.Response[ErrorMethodResponse], r gocharge.Request[string]) error {
		// Use the Error method directly
		_, err := w.Status(gocharge.StatusForbidden).Error("PERMISSION_DENIED", "You do not have permission")
		return err
	})

	response, err := http.Get("http://localhost:8080/testErrorMethod")
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", response.StatusCode)
	}

	var errResp gocharge.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&errResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if errResp.Code != "PERMISSION_DENIED" {
		t.Errorf("Expected error code PERMISSION_DENIED, got %s", errResp.Code)
	}
}

// TestRequestBodyReading verifies request body can be read
func TestRequestBodyReading(t *testing.T) {
	type BodyTestRequest struct {
		Value string `json:"value"`
	}

	type BodyTestResponse struct {
		Received string `json:"received"`
	}

	gocharge.RegisterHandler(server, "/testBody", func(ctx context.Context, w gocharge.Response[BodyTestResponse], r gocharge.Request[BodyTestRequest]) error {
		body, err := r.JSON()
		if err != nil {
			return err
		}

		_, err = w.JSON(BodyTestResponse{Received: body.Value})
		return err
	})

	bodyData := BodyTestRequest{Value: "test-value"}
	bodyBytes, _ := json.Marshal(bodyData)
	response, err := http.Post("http://localhost:8080/testBody", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	body, _ := io.ReadAll(response.Body)
	var resp BodyTestResponse
	json.Unmarshal(body, &resp)

	if resp.Received != "test-value" {
		t.Errorf("Expected 'test-value', got %s", resp.Received)
	}
}
