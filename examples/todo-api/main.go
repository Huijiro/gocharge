package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	gc "github.com/huijiro/go-charge"
	"github.com/huijiro/go-charge/middleware"
)

// ============= Domain Types =============

// Todo represents a single todo item
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// ============= Request/Response Types =============

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TodoResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TodoListResponse struct {
	Todos []TodoResponse `json:"todos"`
	Count int            `json:"count"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ============= In-Memory Storage =============

type TodoStore struct {
	mu    sync.RWMutex
	todos map[string]*Todo
	seq   int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos: make(map[string]*Todo),
		seq:   0,
	}
}

func (s *TodoStore) Create(title string) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	todo := &Todo{
		ID:    fmt.Sprintf("todo-%d", s.seq),
		Title: title,
		Done:  false,
	}
	s.todos[todo.ID] = todo
	return todo
}

func (s *TodoStore) GetAll() []*Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

func (s *TodoStore) GetByID(id string) (*Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	return todo, ok
}

func (s *TodoStore) Update(id string, title string, done bool) (*Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, false
	}

	todo.Title = title
	todo.Done = done
	return todo, true
}

func (s *TodoStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.todos[id]
	if !ok {
		return false
	}

	delete(s.todos, id)
	return true
}

// ============= Handlers =============

// Health check endpoint
func healthHandler(ctx context.Context, w gc.Response[HealthResponse], r gc.Request[string]) error {
	logger := gc.GetLogger()
	logger.Info(ctx, "health check")

	_, err := w.JSON(HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
	return err
}

// Create a new todo
func createTodoHandler(store *TodoStore) gc.HandlerFunc[TodoResponse, CreateTodoRequest] {
	return func(ctx context.Context, w gc.Response[TodoResponse], r gc.Request[CreateTodoRequest]) error {
		logger := gc.GetLogger()

		req, err := r.JSON()
		if err != nil {
			logger.Warn(ctx, "invalid json in create todo request")
			return gc.NewAppError(gc.ErrBadRequest, "Invalid request body")
		}

		// Validate
		if req.Title == "" {
			logger.Warn(ctx, "empty title in create todo request")
			return gc.NewAppError(gc.ErrValidation, "Title is required")
		}

		if len(req.Title) > 100 {
			logger.Warn(ctx, "title too long in create todo request")
			return gc.NewAppError(gc.ErrValidation, "Title must be less than 100 characters")
		}

		// Create todo
		todo := store.Create(req.Title)
		logger.Info(ctx, "todo created", "todo_id", todo.ID, "title", todo.Title)

		_, err = w.Status(gc.StatusCreated).JSON(TodoResponse{
			ID:    todo.ID,
			Title: todo.Title,
			Done:  todo.Done,
		})
		return err
	}
}

// List all todos
func listTodosHandler(store *TodoStore) gc.HandlerFunc[TodoListResponse, string] {
	return func(ctx context.Context, w gc.Response[TodoListResponse], r gc.Request[string]) error {
		logger := gc.GetLogger()
		logger.Info(ctx, "listing todos")

		todos := store.GetAll()
		response := TodoListResponse{
			Todos: make([]TodoResponse, len(todos)),
			Count: len(todos),
		}

		for i, todo := range todos {
			response.Todos[i] = TodoResponse{
				ID:    todo.ID,
				Title: todo.Title,
				Done:  todo.Done,
			}
		}

		_, err := w.JSON(response)
		return err
	}
}

// ============= Main =============

func main() {
	// Setup logger with text output for this example
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewTextHandler(os.Stdout, opts)
	customLogger := gc.NewDefaultLoggerWithHandler(handler)
	gc.SetLogger(customLogger)

	// Create server
	server := gc.New(":8080")
	logger := gc.GetLogger()

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

	// Start server
	logger.Info(context.Background(), "server starting", "addr", ":8080")
	if err := server.ListenAndServe(); err != nil {
		logger.Error(context.Background(), "server error", err)
		os.Exit(1)
	}
}
