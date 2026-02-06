package gocharge_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

type ItemResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestPathParamString(t *testing.T) {
	gocharge.RegisterHandler(server, "/items/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().String("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{
			ID:   id,
			Name: "Item " + id,
		}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items/test-123")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var resp ItemResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if resp.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", resp.ID)
	}
}

func TestPathParamInt(t *testing.T) {
	gocharge.RegisterHandler(server, "/items-int/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().Int("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{
			ID:    "",
			Name:  "Item",
			Count: id,
		}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items-int/42")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var resp ItemResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if resp.Count != 42 {
		t.Errorf("Expected count 42, got %d", resp.Count)
	}
}

func TestPathParamIntInvalid(t *testing.T) {
	gocharge.RegisterHandler(server, "/items-int-invalid/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().Int("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{Count: id}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items-int-invalid/not-a-number")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
	}

	var errResp map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&errResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	errMsg, ok := errResp["message"].(string)
	if !ok {
		t.Errorf("Expected error message in response")
	}

	if !strings.Contains(errMsg, "invalid id format: expected integer, got 'not-a-number'") {
		t.Errorf("Expected error message to contain 'invalid id format: expected integer, got 'not-a-number'', got: %s", errMsg)
	}
}

func TestPathParamInt64(t *testing.T) {
	gocharge.RegisterHandler(server, "/items-int64/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().Int64("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{
			ID:    "",
			Name:  "Item",
			Count: int(id),
		}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items-int64/9223372036854775807")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var resp ItemResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if resp.Count != 9223372036854775807 {
		t.Errorf("Expected count 9223372036854775807, got %d", resp.Count)
	}
}

func TestPathParamInt64Invalid(t *testing.T) {
	gocharge.RegisterHandler(server, "/items-int64-invalid/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().Int64("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{Count: int(id)}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items-int64-invalid/abc")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
	}

	var errResp map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&errResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	errMsg, ok := errResp["message"].(string)
	if !ok {
		t.Errorf("Expected error message in response")
	}

	if !strings.Contains(errMsg, "invalid id format: expected integer, got 'abc'") {
		t.Errorf("Expected error message to contain 'invalid id format: expected integer, got 'abc'', got: %s", errMsg)
	}
}

func TestPathParamUUID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	gocharge.RegisterHandler(server, "/items-uuid/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().UUID("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{
			ID:   id,
			Name: "Item",
		}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items-uuid/" + validUUID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var resp ItemResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if resp.ID != validUUID {
		t.Errorf("Expected ID %s, got %s", validUUID, resp.ID)
	}
}

func TestPathParamUUIDInvalidFormat(t *testing.T) {
	invalidUUID := "not-a-uuid"

	gocharge.RegisterHandler(server, "/items-uuid-invalid/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().UUID("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{ID: id}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items-uuid-invalid/" + invalidUUID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
	}

	var errResp map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&errResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	errMsg, ok := errResp["message"].(string)
	if !ok {
		t.Errorf("Expected error message in response")
	}

	if !strings.Contains(errMsg, "invalid id format: expected UUID") {
		t.Errorf("Expected error message to contain 'invalid id format: expected UUID', got: %s", errMsg)
	}
}

func TestPathParamUUIDInvalidHex(t *testing.T) {
	invalidUUID := "550e8400-e29b-41d4-a716-44665544000g" // 'g' is not hex

	gocharge.RegisterHandler(server, "/items-uuid-hex/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		id, err := r.PathParam().UUID("id")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{ID: id}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/items-uuid-hex/" + invalidUUID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
	}

	var errResp map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&errResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	errMsg, ok := errResp["message"].(string)
	if !ok {
		t.Errorf("Expected error message in response")
	}

	if !strings.Contains(errMsg, "invalid id format: expected UUID") {
		t.Errorf("Expected error message to contain 'invalid id format: expected UUID', got: %s", errMsg)
	}
}

func TestPathParamMultiple(t *testing.T) {
	gocharge.RegisterHandler(server, "/org/{orgId}/projects/{projectId}", func(ctx context.Context, w gocharge.Response[map[string]interface{}], r gocharge.Request[string]) error {
		orgID, err := r.PathParam().String("orgId")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, "orgId: "+err.Error())
		}

		projectID, err := r.PathParam().Int("projectId")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, "projectId: "+err.Error())
		}

		resp := map[string]interface{}{
			"orgId":     orgID,
			"projectId": projectID,
		}
		_, err = w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/org/acme/projects/123")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var resp map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if resp["orgId"] != "acme" {
		t.Errorf("Expected orgId 'acme', got '%v'", resp["orgId"])
	}

	if resp["projectId"] != float64(123) {
		t.Errorf("Expected projectId 123, got '%v'", resp["projectId"])
	}
}

func TestPathParamMissing(t *testing.T) {
	gocharge.RegisterHandler(server, "/items-missing/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		// Try to get a parameter that doesn't exist in the route
		id, err := r.PathParam().String("nonexistent")
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, err.Error())
		}

		resp := ItemResponse{ID: id}
		_, err = w.JSON(resp)
		return err
	})

	// This request should match the route
	response, err := http.Get("http://localhost:8080/items-missing/123")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	// Should get a 400 error because the param doesn't exist
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
	}

	var errResp map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&errResp)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	errMsg, ok := errResp["message"].(string)
	if !ok {
		t.Errorf("Expected error message in response")
	}

	if !strings.Contains(errMsg, "missing path parameter: nonexistent") {
		t.Errorf("Expected error message to contain 'missing path parameter: nonexistent', got: %s", errMsg)
	}
}
