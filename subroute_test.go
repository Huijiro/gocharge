package gocharge_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

// Test data structures for sub-route tests
type UserCreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListResponse struct {
	Items []UserResponse `json:"items"`
	Count int            `json:"count"`
}

// TestSubRouteBasicPattern demonstrates basic sub-route structure
// Pattern: /api/users, /api/users/{id}, etc.
func TestSubRouteBasicPattern(t *testing.T) {
	// Base route
	gocharge.RegisterHandler(server, "/api/users", func(ctx context.Context, w gocharge.Response[ListResponse], r gocharge.Request[string]) error {
		resp := ListResponse{
			Items: []UserResponse{
				{ID: "1", Name: "Alice", Email: "alice@example.com"},
				{ID: "2", Name: "Bob", Email: "bob@example.com"},
			},
			Count: 2,
		}
		_, err := w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/api/users")
	if err != nil {
		t.Errorf("Error making request: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	var listResp ListResponse
	json.NewDecoder(response.Body).Decode(&listResp)

	if listResp.Count != 2 {
		t.Errorf("Expected 2 items, got %d", listResp.Count)
	}
}

// TestSubRouteWithPrefix demonstrates organizing handlers with path prefixes
// Pattern: All user endpoints start with /api/v1/users
func TestSubRouteWithPrefix(t *testing.T) {
	prefix := "/api/v1"

	// Create user
	gocharge.RegisterHandler(server, prefix+"/users", func(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[UserCreateRequest]) error {
		req, err := r.JSON()
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, "Invalid request body")
		}

		resp := UserResponse{
			ID:    "new-user",
			Name:  req.Name,
			Email: req.Email,
		}
		_, err = w.Status(gocharge.StatusCreated).JSON(resp)
		return err
	})

	// Test create user
	body := UserCreateRequest{Name: "Charlie", Email: "charlie@example.com"}
	bodyBytes, _ := json.Marshal(body)
	response, err := http.Post("http://localhost:8080/api/v1/users", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Errorf("Error making request: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", response.StatusCode)
	}

	var userResp UserResponse
	json.NewDecoder(response.Body).Decode(&userResp)

	if userResp.Name != "Charlie" {
		t.Errorf("Expected name Charlie, got %s", userResp.Name)
	}
}

// TestSubRouteHierarchy demonstrates nested resource patterns
// Pattern: /api/users/{userId}/posts, /api/users/{userId}/settings
func TestSubRouteHierarchy(t *testing.T) {
	type PostResponse struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	type UserPostsResponse struct {
		UserID string         `json:"user_id"`
		Posts  []PostResponse `json:"posts"`
	}

	// Nested resource: user posts
	gocharge.RegisterHandler(server, "/api/users/{userId}/posts", func(ctx context.Context, w gocharge.Response[UserPostsResponse], r gocharge.Request[string]) error {
		// In a real app, you'd extract userId from the URL path
		resp := UserPostsResponse{
			UserID: "user-123",
			Posts: []PostResponse{
				{ID: "post-1", Title: "First Post", Body: "Content..."},
				{ID: "post-2", Title: "Second Post", Body: "More content..."},
			},
		}
		_, err := w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/api/users/{userId}/posts")
	if err != nil {
		t.Errorf("Error making request: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	var postsResp UserPostsResponse
	json.NewDecoder(response.Body).Decode(&postsResp)

	if len(postsResp.Posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(postsResp.Posts))
	}
}

// TestSubRouteMultipleLevels demonstrates multiple levels of nesting
// Pattern: /api/v1/organizations/{orgId}/projects/{projectId}/tasks
func TestSubRouteMultipleLevels(t *testing.T) {
	type TaskResponse struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	type TaskListResponse struct {
		Tasks []TaskResponse `json:"tasks"`
	}

	// Deeply nested resource
	gocharge.RegisterHandler(server, "/api/v1/organizations/{orgId}/projects/{projectId}/tasks", func(ctx context.Context, w gocharge.Response[TaskListResponse], r gocharge.Request[string]) error {
		resp := TaskListResponse{
			Tasks: []TaskResponse{
				{ID: "task-1", Title: "Task One"},
				{ID: "task-2", Title: "Task Two"},
			},
		}
		_, err := w.JSON(resp)
		return err
	})

	response, err := http.Get("http://localhost:8080/api/v1/organizations/{orgId}/projects/{projectId}/tasks")
	if err != nil {
		t.Errorf("Error making request: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	var tasksResp TaskListResponse
	json.NewDecoder(response.Body).Decode(&tasksResp)

	if len(tasksResp.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasksResp.Tasks))
	}
}

// TestSubRouteSeparationOfConcerns demonstrates organizing handlers by domain
// Pattern: Different handler groups for different resources
func TestSubRouteSeparationOfConcerns(t *testing.T) {
	type ProfileResponse struct {
		UserID string `json:"user_id"`
		Bio    string `json:"bio"`
	}

	type SettingsResponse struct {
		UserID        string `json:"user_id"`
		Theme         string `json:"theme"`
		Notifications bool   `json:"notifications"`
	}

	// Profile endpoint
	gocharge.RegisterHandler(server, "/api/profile", func(ctx context.Context, w gocharge.Response[ProfileResponse], r gocharge.Request[string]) error {
		_, err := w.JSON(ProfileResponse{UserID: "user-1", Bio: "Bio content"})
		return err
	})

	// Settings endpoint
	gocharge.RegisterHandler(server, "/api/settings", func(ctx context.Context, w gocharge.Response[SettingsResponse], r gocharge.Request[string]) error {
		_, err := w.JSON(SettingsResponse{UserID: "user-1", Theme: "dark", Notifications: true})
		return err
	})

	// Test profile
	resp1, _ := http.Get("http://localhost:8080/api/profile")
	defer resp1.Body.Close()

	var profile ProfileResponse
	json.NewDecoder(resp1.Body).Decode(&profile)

	if profile.Bio != "Bio content" {
		t.Errorf("Expected bio content, got %s", profile.Bio)
	}

	// Test settings
	resp2, _ := http.Get("http://localhost:8080/api/settings")
	defer resp2.Body.Close()

	var settings SettingsResponse
	json.NewDecoder(resp2.Body).Decode(&settings)

	if settings.Theme != "dark" {
		t.Errorf("Expected theme dark, got %s", settings.Theme)
	}
}

// TestSubRouteErrorHandlingInSubRoutes demonstrates error handling in sub-routes
func TestSubRouteErrorHandlingInSubRoutes(t *testing.T) {
	type GetUserRequest struct {
		ID string `json:"id"`
	}

	// Sub-route that may return errors
	gocharge.RegisterHandler(server, "/api/v2/users/{id}", func(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[GetUserRequest]) error {
		req, err := r.JSON()
		if err != nil {
			return gocharge.NewAppError(gocharge.ErrBadRequest, "Invalid request")
		}

		if req.ID == "" {
			return gocharge.NewAppError(gocharge.ErrValidation, "User ID is required")
		}

		if req.ID == "999" {
			return gocharge.NewAppError(gocharge.ErrNotFound, "User not found")
		}

		resp := UserResponse{ID: req.ID, Name: "Found User", Email: "user@example.com"}
		_, err = w.JSON(resp)
		return err
	})

	// Test with valid ID
	body := GetUserRequest{ID: "123"}
	bodyBytes, _ := json.Marshal(body)
	response, _ := http.Post("http://localhost:8080/api/v2/users/{id}", "application/json", bytes.NewBuffer(bodyBytes))
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	// Test with non-existent ID
	body2 := GetUserRequest{ID: "999"}
	bodyBytes2, _ := json.Marshal(body2)
	response2, _ := http.Post("http://localhost:8080/api/v2/users/{id}", "application/json", bytes.NewBuffer(bodyBytes2))
	defer response2.Body.Close()

	if response2.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", response2.StatusCode)
	}

	var errResp gocharge.ErrorResponse
	json.NewDecoder(response2.Body).Decode(&errResp)

	if errResp.Code != string(gocharge.ErrNotFound) {
		t.Errorf("Expected error code NOT_FOUND, got %s", errResp.Code)
	}
}

// TestSubRouteRESTfulPattern demonstrates RESTful API pattern with sub-routes
// Pattern: GET /items (list), POST /items (create), GET /items/{id} (get), PUT /items/{id} (update)
func TestSubRouteRESTfulPattern(t *testing.T) {
	type ItemResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type ItemListResponse struct {
		Items []ItemResponse `json:"items"`
	}

	// List items
	gocharge.RegisterHandler(server, "/rest/items", func(ctx context.Context, w gocharge.Response[ItemListResponse], r gocharge.Request[string]) error {
		resp := ItemListResponse{
			Items: []ItemResponse{
				{ID: "1", Name: "Item 1"},
				{ID: "2", Name: "Item 2"},
			},
		}
		_, err := w.JSON(resp)
		return err
	})

	// Get single item (in real app, would extract ID from path)
	gocharge.RegisterHandler(server, "/rest/items/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
		resp := ItemResponse{ID: "123", Name: "Specific Item"}
		_, err := w.JSON(resp)
		return err
	})

	// Test list
	response, _ := http.Get("http://localhost:8080/rest/items")
	defer response.Body.Close()

	var listResp ItemListResponse
	json.NewDecoder(response.Body).Decode(&listResp)

	if len(listResp.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(listResp.Items))
	}

	// Test get single
	response2, _ := http.Get("http://localhost:8080/rest/items/{id}")
	defer response2.Body.Close()

	var itemResp ItemResponse
	json.NewDecoder(response2.Body).Decode(&itemResp)

	if itemResp.ID != "123" {
		t.Errorf("Expected ID 123, got %s", itemResp.ID)
	}
}

// TestSubRouteWithSharedMiddleware demonstrates organizing multiple sub-routes
// that can share middleware in real applications
func TestSubRouteWithSharedMiddleware(t *testing.T) {
	type AuthResponse struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}

	// Sub-routes that would benefit from shared middleware
	// In a real app, these would all use the same auth middleware
	publicHandler := func(ctx context.Context, w gocharge.Response[AuthResponse], r gocharge.Request[string]) error {
		resp := AuthResponse{
			Message: "Public endpoint",
			Success: true,
		}
		_, err := w.JSON(resp)
		return err
	}

	gocharge.RegisterHandler(server, "/public/secure/data", publicHandler)
	gocharge.RegisterHandler(server, "/public/secure/profile", publicHandler)

	// Test protected routes
	response, _ := http.Get("http://localhost:8080/public/secure/data")
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	var authResp AuthResponse
	json.NewDecoder(response.Body).Decode(&authResp)

	if !authResp.Success {
		t.Error("Expected success to be true")
	}

	// Test second route
	response2, _ := http.Get("http://localhost:8080/public/secure/profile")
	defer response2.Body.Close()

	if response2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response2.StatusCode)
	}
}

// TestSubRoutePathVariations demonstrates different path patterns
func TestSubRoutePathVariations(t *testing.T) {
	type SimpleResponse struct {
		Path string `json:"path"`
	}

	paths := map[string]string{
		"/admin/dashboard":  "admin/dashboard",
		"/api/v1/health":    "api/v1/health",
		"/public/assets":    "public/assets",
		"/internal/metrics": "internal/metrics",
		"/webhooks/github":  "webhooks/github",
		"/graphql":          "graphql",
	}

	for path, desc := range paths {
		gocharge.RegisterHandler(server, path, func(ctx context.Context, w gocharge.Response[SimpleResponse], r gocharge.Request[string]) error {
			_, err := w.JSON(SimpleResponse{Path: desc})
			return err
		})
	}

	// Test one of them
	response, _ := http.Get("http://localhost:8080/admin/dashboard")
	defer response.Body.Close()

	var resp SimpleResponse
	json.NewDecoder(response.Body).Decode(&resp)

	if resp.Path != "admin/dashboard" {
		t.Errorf("Expected path admin/dashboard, got %s", resp.Path)
	}
}
