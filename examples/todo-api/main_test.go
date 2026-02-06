package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	gc "github.com/huijiro/go-charge"
	"github.com/huijiro/go-charge/middleware"
)

// startTestServer creates and starts a test server on a random available port
func startTestServer(t *testing.T) (*gc.Server, string) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	addr := listener.Addr().String()

	// Create server with the listener's address
	server := gc.New(addr)

	// Setup middleware
	chain := server.Middleware()
	chain.Use(middleware.LoggingMiddleware)
	chain.Use(middleware.RecoveryMiddleware)

	// Create data store
	store := NewTodoStore()

	// Register handlers
	gc.RegisterHandler(server, "/health", healthHandler)
	gc.RegisterHandler(server, "POST /api/todos", createTodoHandler(store))
	gc.RegisterHandler(server, "GET /api/todos", listTodosHandler(store))
	gc.RegisterHandler(server, "GET /api/todos/{id}", GetTodoByIDHandler(store))
	gc.RegisterHandler(server, "PUT /api/todos/{id}", UpdateTodoHandler(store))
	gc.RegisterHandler(server, "DELETE /api/todos/{id}", DeleteTodoHandler(store))

	// Start server in a goroutine
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	return server, "http://" + addr
}

func TestTodoAPIFullWorkflow(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	// 1. Create first todo
	createReq := CreateTodoRequest{Title: "Learn GoCharge"}
	createBody, _ := json.Marshal(createReq)
	resp, err := http.Post(baseURL+"/api/todos", "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var todo1 TodoResponse
	json.NewDecoder(resp.Body).Decode(&todo1)
	if todo1.Title != "Learn GoCharge" || todo1.Done {
		t.Errorf("Unexpected todo response: %+v", todo1)
	}
	todo1ID := todo1.ID

	// 2. Create second todo
	createReq2 := CreateTodoRequest{Title: "Build API"}
	createBody2, _ := json.Marshal(createReq2)
	resp2, _ := http.Post(baseURL+"/api/todos", "application/json", bytes.NewBuffer(createBody2))
	defer resp2.Body.Close()

	var todo2 TodoResponse
	json.NewDecoder(resp2.Body).Decode(&todo2)
	todo2ID := todo2.ID

	// 3. List todos
	resp3, _ := http.Get(baseURL + "/api/todos")
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp3.StatusCode)
	}

	var listResp TodoListResponse
	json.NewDecoder(resp3.Body).Decode(&listResp)
	if listResp.Count != 2 {
		t.Errorf("Expected 2 todos, got %d", listResp.Count)
	}

	// 4. Get single todo
	resp4, _ := http.Get(baseURL + "/api/todos/" + todo1ID)
	defer resp4.Body.Close()

	if resp4.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp4.StatusCode)
	}

	var getTodo TodoResponse
	json.NewDecoder(resp4.Body).Decode(&getTodo)
	if getTodo.ID != todo1ID || getTodo.Title != "Learn GoCharge" {
		t.Errorf("Unexpected todo: %+v", getTodo)
	}

	// 5. Update todo
	updateReq := UpdateTodoRequest{Title: "Learn GoCharge Framework", Done: true}
	updateBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", baseURL+"/api/todos/"+todo1ID, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp5, _ := http.DefaultClient.Do(req)
	defer resp5.Body.Close()

	if resp5.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp5.StatusCode)
	}

	var updatedTodo TodoResponse
	json.NewDecoder(resp5.Body).Decode(&updatedTodo)
	if updatedTodo.Title != "Learn GoCharge Framework" || !updatedTodo.Done {
		t.Errorf("Unexpected updated todo: %+v", updatedTodo)
	}

	// 6. Delete todo
	req2, _ := http.NewRequest("DELETE", baseURL+"/api/todos/"+todo2ID, nil)
	resp6, _ := http.DefaultClient.Do(req2)
	defer resp6.Body.Close()

	if resp6.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, resp6.StatusCode)
	}

	// 7. Verify todo was deleted
	resp7, _ := http.Get(baseURL + "/api/todos/" + todo2ID)
	defer resp7.Body.Close()

	if resp7.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp7.StatusCode)
	}

	// 8. Verify list now has only 1 todo
	resp8, _ := http.Get(baseURL + "/api/todos")
	defer resp8.Body.Close()

	var finalList TodoListResponse
	json.NewDecoder(resp8.Body).Decode(&finalList)
	if finalList.Count != 1 {
		t.Errorf("Expected 1 todo after deletion, got %d", finalList.Count)
	}
}

func TestCreateTodoValidation(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		title          string
		expectedStatus int
		expectedError  string
	}{
		{"Empty title", "", http.StatusBadRequest, "Title is required"},
		{"Valid title", "Test Todo", http.StatusCreated, ""},
		{"Title too long", string(make([]byte, 101)), http.StatusBadRequest, "Title must be less than 100 characters"},
		{"Long valid title", string(make([]byte, 100)), http.StatusCreated, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createReq := CreateTodoRequest{Title: test.title}
			createBody, _ := json.Marshal(createReq)
			resp, _ := http.Post(baseURL+"/api/todos", "application/json", bytes.NewBuffer(createBody))
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, resp.StatusCode)
			}

			if test.expectedError != "" {
				var errResp map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&errResp)
				errMsg, ok := errResp["message"].(string)
				if !ok || errMsg != test.expectedError {
					t.Errorf("Expected error '%s', got '%v'", test.expectedError, errResp)
				}
			}
		})
	}
}

func TestGetTodoNotFound(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	resp, _ := http.Get(baseURL + "/api/todos/nonexistent")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	var errResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	errMsg, ok := errResp["message"].(string)
	if !ok || errMsg != "Todo not found" {
		t.Errorf("Expected 'Todo not found' error, got '%v'", errResp)
	}
}

func TestUpdateTodoValidation(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	// Create a todo first
	createReq := CreateTodoRequest{Title: "Original"}
	createBody, _ := json.Marshal(createReq)
	resp, _ := http.Post(baseURL+"/api/todos", "application/json", bytes.NewBuffer(createBody))
	defer resp.Body.Close()

	var todo TodoResponse
	json.NewDecoder(resp.Body).Decode(&todo)
	todoID := todo.ID

	// Test update with empty title
	updateReq := UpdateTodoRequest{Title: "", Done: false}
	updateBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", baseURL+"/api/todos/"+todoID, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp2, _ := http.DefaultClient.Do(req)
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp2.StatusCode)
	}

	// Test update with valid data
	updateReq2 := UpdateTodoRequest{Title: "Updated", Done: true}
	updateBody2, _ := json.Marshal(updateReq2)
	req2, _ := http.NewRequest("PUT", baseURL+"/api/todos/"+todoID, bytes.NewBuffer(updateBody2))
	req2.Header.Set("Content-Type", "application/json")
	resp3, _ := http.DefaultClient.Do(req2)
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp3.StatusCode)
	}
}

func TestUpdateTodoNotFound(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	updateReq := UpdateTodoRequest{Title: "Update", Done: false}
	updateBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", baseURL+"/api/todos/nonexistent", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestDeleteTodoNotFound(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	req, _ := http.NewRequest("DELETE", baseURL+"/api/todos/nonexistent", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestHealthEndpoint(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	resp, _ := http.Get(baseURL + "/health")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var health HealthResponse
	json.NewDecoder(resp.Body).Decode(&health)
	if health.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", health.Status)
	}
}

func TestPathParamExtraction(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	// Create a todo with ID like "todo-1"
	createReq := CreateTodoRequest{Title: "Path Param Test"}
	createBody, _ := json.Marshal(createReq)
	resp, _ := http.Post(baseURL+"/api/todos", "application/json", bytes.NewBuffer(createBody))
	defer resp.Body.Close()

	var todo TodoResponse
	json.NewDecoder(resp.Body).Decode(&todo)

	// The ID should be extractable with string path param
	resp2, _ := http.Get(baseURL + "/api/todos/" + todo.ID)
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp2.StatusCode)
	}

	var retrieved TodoResponse
	json.NewDecoder(resp2.Body).Decode(&retrieved)
	if retrieved.ID != todo.ID {
		t.Errorf("Expected ID '%s', got '%s'", todo.ID, retrieved.ID)
	}
}

func TestInvalidJSON(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	// Send invalid JSON
	resp, _ := http.Post(baseURL+"/api/todos", "application/json", bytes.NewBuffer([]byte("invalid json")))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var errResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&errResp)
	if errResp["message"] == nil {
		t.Errorf("Expected error message in response")
	}
}

func TestEmptyListInitially(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	resp, _ := http.Get(baseURL + "/api/todos")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var listResp TodoListResponse
	json.NewDecoder(resp.Body).Decode(&listResp)
	if listResp.Count != 0 {
		t.Errorf("Expected 0 todos initially, got %d", listResp.Count)
	}
}

func TestMultipleTodosOrderIndependent(t *testing.T) {
	server, baseURL := startTestServer(t)
	defer server.Close()

	// Create multiple todos with different titles
	titles := []string{"First", "Second", "Third"}
	ids := make([]string, 0)

	for i, title := range titles {
		createReq := CreateTodoRequest{Title: title}
		createBody, _ := json.Marshal(createReq)
		resp, _ := http.Post(baseURL+"/api/todos", "application/json", bytes.NewBuffer(createBody))
		defer resp.Body.Close()

		var todo TodoResponse
		json.NewDecoder(resp.Body).Decode(&todo)
		ids = append(ids, todo.ID)

		if todo.Title != title {
			t.Errorf("Todo %d: expected title '%s', got '%s'", i, title, todo.Title)
		}
	}

	// Verify we can retrieve all of them by ID
	for i, id := range ids {
		resp, _ := http.Get(baseURL + "/api/todos/" + id)
		defer resp.Body.Close()

		var todo TodoResponse
		json.NewDecoder(resp.Body).Decode(&todo)
		if todo.Title != titles[i] {
			t.Errorf("Retrieved todo %d: expected title '%s', got '%s'", i, titles[i], todo.Title)
		}
	}
}
