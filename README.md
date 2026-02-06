# GoCharge

A **production-ready, fully typed HTTP backend framework** for Go 1.18+ using only the Go standard library.

[![Go Version](https://img.shields.io/badge/go-1.25.6+-00ADD8?style=flat-square&logo=go)](https://golang.org/doc/devel/release)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE.md)
[![Tests](https://img.shields.io/badge/tests-54+-green?style=flat-square)](.github/workflows/go.yml)
[![Zero Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen?style=flat-square)](#no-external-dependencies)

## Why GoCharge?

GoCharge eliminates boilerplate and provides **compile-time type safety** for your HTTP handlers without external dependencies. Requires Go 1.25.6+ to take advantage of all features including method-aware routing. Build fast, maintainable APIs with confidence.

```go
// Type-safe handlers with automatic error handling
func createUserHandler(ctx context.Context, w gocharge.Response[User], r gocharge.Request[CreateUserRequest]) error {
    req, err := r.JSON()
    if err != nil {
        return gocharge.NewAppError(gocharge.ErrBadRequest, "Invalid JSON")
    }
    
    // Business logic...
    user := createUser(req)
    
    _, err = w.Status(gocharge.StatusCreated).JSON(user)
    return err
}

// Register with type-safe path parameters
gc.RegisterHandler(server, "POST /api/users", createUserHandler)
gc.RegisterHandler(server, "GET /api/users/{id}", getUserHandler)
gc.RegisterHandler(server, "PUT /api/users/{id}", updateUserHandler)
gc.RegisterHandler(server, "DELETE /api/users/{id}", deleteUserHandler)
```

## Key Features

### ✨ Type-Safe Everything

- **Typed Handlers**: Generic request/response types with compile-time safety
- **Typed Errors**: Strongly-typed error codes with automatic HTTP status mapping
- **Typed Context**: Type-safe context value storage and retrieval
- **Path Parameters**: Type-safe extraction with built-in converters (string, int, int64, UUID)

### 🚀 Built for Production

- **Zero External Dependencies**: Only Go stdlib (except tests)
- **Structured Logging**: Built-in slog integration with custom logger support
- **Error Handling**: Automatic error encoding with customizable error responses
- **Middleware System**: Composable middleware chain with built-in logging and recovery
- **Method-Aware Routes**: Go 1.22+ native HTTP method support in paths

### 📦 Developer Experience

- **Minimal Boilerplate**: No interfaces to implement, just write handlers
- **Response Chaining**: Fluent API for building responses
- **Automatic JSON**: Built-in JSON encoding/decoding
- **Request Validation**: Type-safe request parsing with error propagation
- **Comprehensive Documentation**: Inline examples for every feature

### 🧪 Well-Tested

- **54+ Tests**: Comprehensive test coverage
- **Zero Linting Issues**: golangci-lint clean
- **Security Scanned**: gosec validation
- **Multi-Version Testing**: Tested on Go 1.21, 1.22, 1.25

## Installation

```bash
go get github.com/huijiro/go-charge
```

**Requirements:**
- Go 1.25.6 or later (for method-aware paths and full feature support)
- No external dependencies

## Quick Start

### 1. Create a Server

```go
package main

import (
    "context"
    gc "github.com/huijiro/go-charge"
)

func main() {
    server := gc.New(":8080")
    
    // Register handlers
    gc.RegisterHandler(server, "GET /health", healthHandler)
    gc.RegisterHandler(server, "POST /api/todos", createTodoHandler)
    
    // Start server
    server.ListenAndServe()
}
```

### 2. Define Request/Response Types

```go
type CreateTodoRequest struct {
    Title string `json:"title"`
}

type TodoResponse struct {
    ID    string `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}
```

### 3. Write Type-Safe Handlers

```go
func createTodoHandler(ctx context.Context, w gc.Response[TodoResponse], r gc.Request[CreateTodoRequest]) error {
    logger := gc.GetLogger()
    
    // Parse and validate request
    req, err := r.JSON()
    if err != nil {
        logger.Warn(ctx, "invalid json")
        return gc.NewAppError(gc.ErrBadRequest, "Invalid request body")
    }
    
    if req.Title == "" {
        return gc.NewAppError(gc.ErrValidation, "Title is required")
    }
    
    // Business logic
    todo := createTodo(req.Title)
    logger.Info(ctx, "todo created", "id", todo.ID)
    
    // Send response
    _, err = w.Status(gc.StatusCreated).JSON(TodoResponse{
        ID:    todo.ID,
        Title: todo.Title,
        Done:  todo.Done,
    })
    return err
}
```

## Core Concepts

### Request/Response Types

```go
// Request[T] wraps http.Request with typed body
type Request[T any] struct {
    Data T
    http.Request
}

// Response[T] wraps http.ResponseWriter with typed JSON
type Response[T any] struct {
    ResponseWriter http.ResponseWriter
    Data           T
    StatusCode     int
}
```

### Type-Safe Errors

```go
// Create typed errors
err := gc.NewAppError(gc.ErrBadRequest, "Invalid email")
err := gc.NewAppError(gc.ErrNotFound, "User not found")
err := gc.NewAppError(gc.ErrValidation, "Password too short")

// Built-in error codes with automatic HTTP status mapping
gc.ErrBadRequest        // 400
gc.ErrUnauthorized      // 401
gc.ErrForbidden         // 403
gc.ErrNotFound          // 404
gc.ErrValidation        // 400
gc.ErrInternalError     // 500
```

### Path Parameter Extraction

```go
func getUserHandler(ctx context.Context, w gc.Response[User], r gc.Request[string]) error {
    // Extract string parameter
    userID, err := r.PathParam().String("id")
    if err != nil {
        return gc.NewAppError(gc.ErrBadRequest, err.Error())
    }
    
    user, found := getUser(userID)
    if !found {
        return gc.NewAppError(gc.ErrNotFound, "User not found")
    }
    
    _, err = w.JSON(user)
    return err
}

// Supported parameter types:
id, err := r.PathParam().String("id")     // Extract string
count, err := r.PathParam().Int("count")  // Extract int
offset, err := r.PathParam().Int64("offset") // Extract int64
uuid, err := r.PathParam().UUID("id")     // Extract & validate RFC 4122 UUID

// Register with method-aware paths
gc.RegisterHandler(server, "GET /api/users/{id}", getUserHandler)
gc.RegisterHandler(server, "PUT /api/users/{id}", updateUserHandler)
gc.RegisterHandler(server, "DELETE /api/users/{id}", deleteUserHandler)
```

### Structured Logging

```go
logger := gc.GetLogger()

// Log with context and structured fields
logger.Info(ctx, "user created", "user_id", user.ID, "email", user.Email)
logger.Warn(ctx, "invalid input", "reason", "email_taken")
logger.Error(ctx, "database error", err)

// Custom logger support
type CustomLogger struct{}
func (l *CustomLogger) Info(ctx context.Context, msg string, fields ...interface{}) { }
// Implement Logger interface...
gc.SetLogger(customLogger)
```

### Middleware

```go
// Access built-in middleware
chain := server.Middleware()
chain.Use(middleware.LoggingMiddleware)  // Logs all requests/responses
chain.Use(middleware.RecoveryMiddleware) // Catches panics, returns 500

// Create custom middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }
        // Verify token and add to context...
        next.ServeHTTP(w, r)
    })
}

chain.Use(authMiddleware)
```

### Response Chaining

```go
// Fluent API for building responses
_, err = w
    .Status(gc.StatusCreated)
    .JSON(TodoResponse{
        ID:    todo.ID,
        Title: todo.Title,
        Done:  todo.Done,
    })
return err

// Or write raw content
_, err = w
    .Status(gc.StatusOK)
    .Write([]byte("Hello, World!"))
return err
```

### Type-Safe Context Values

```go
// Store values in request context
ctx = gc.WithUser(ctx, user)
ctx = gc.WithTraceID(ctx, traceID)

// Retrieve values with type safety
user, hasUser := gc.UserFromContext(ctx)
traceID, hasTraceID := gc.TraceIDFromContext(ctx)

// Custom typed values
ctx = gc.WithValue(ctx, "key", value)
val, hasVal := gc.Value[MyType](ctx, "key")
```

## Examples

### Todo API

A complete, production-ready TODO API example demonstrating all GoCharge features:

```bash
cd examples/todo-api
go run main.go
```

**Endpoints:**
- `POST /api/todos` - Create todo
- `GET /api/todos` - List todos
- `GET /api/todos/{id}` - Get single todo
- `PUT /api/todos/{id}` - Update todo
- `DELETE /api/todos/{id}` - Delete todo

**Features:**
- Type-safe handlers
- Path parameter extraction
- Request validation
- Structured logging
- Error handling
- Middleware integration
- Integration tests

See [`examples/todo-api/README.md`](examples/todo-api/README.md) for detailed documentation and curl examples.

### Hello World

A minimal example showing the basics:

```bash
cd examples/hello-world
go run main.go
```

See [`examples/hello-world`](examples/hello-world) for the complete example.

## Architecture

### Package Structure

```
gocharge/
├── gocharge.go              # Server initialization and core types
├── handler.go               # Handler registration and wrapping
├── request.go               # Request wrapper with JSON decoding & path params
├── response.go              # Response wrapper with chaining
├── errors.go                # Error types and encoding
├── context.go               # Context helpers and typed values
├── logger.go                # Logger interface and default implementation
├── middleware.go            # Middleware chain system
└── middleware/
    ├── logger.go            # Built-in logging middleware
    └── recovery.go          # Built-in panic recovery middleware
```

### Request/Response Flow

```
1. HTTP Request
   ↓
2. RegisterHandler wraps with type safety
   ↓
3. Extract path parameters (Go 1.22+ PathValue)
   ↓
4. Create typed Request[T] with body data
   ↓
5. Create typed Response[W]
   ↓
6. Call user handler with context
   ↓
7. Handler returns error or nil
   ↓
8. Automatic error encoding (if error returned)
   ↓
9. HTTP Response
```

## Error Handling

GoCharge provides automatic error handling with customizable encoders:

```go
// Default error encoder produces:
{
    "code": "BAD_REQUEST",
    "message": "Invalid email format"
}

// Customize error responses
type CustomErrorEncoder struct{}
func (e *CustomErrorEncoder) Encode(ctx context.Context, err error) (interface{}, int) {
    // Custom error formatting logic
    return response, statusCode
}
gc.SetErrorEncoder(server, customEncoder)
```

**Error Response Example:**
```json
{
    "code": "VALIDATION_ERROR",
    "message": "Title must be between 1 and 100 characters"
}
```

## Testing

GoCharge applications are straightforward to test:

```go
func TestCreateTodo(t *testing.T) {
    server := gc.New(":8081")
    gc.RegisterHandler(server, "POST /api/todos", createTodoHandler)
    
    // Start server
    go server.ListenAndServe()
    defer server.Close()
    
    // Make request
    resp, err := http.Post(
        "http://localhost:8081/api/todos",
        "application/json",
        bytes.NewBuffer([]byte(`{"title":"Test"}`)),
    )
    
    // Assert response
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

See [`examples/todo-api/main_test.go`](examples/todo-api/main_test.go) for comprehensive integration tests.

## Performance

GoCharge is built for performance:

- **No reflection overhead**: Types resolved at compile time
- **Minimal allocations**: Efficient request/response wrapping
- **Direct stdlib integration**: No abstraction layers
- **Zero dependencies**: Direct http.ServeMux usage

## Best Practices

### 1. Separate Concerns

```go
// Handlers in handlers.go
// Models in models.go
// Business logic in services/
// Middleware in middleware/
```

### 2. Type Your Requests and Responses

```go
// Always define explicit types
type CreateUserRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}
```

### 3. Validate Early

```go
req, err := r.JSON()
if err != nil {
    return gc.NewAppError(gc.ErrBadRequest, "Invalid JSON")
}

if req.Email == "" {
    return gc.NewAppError(gc.ErrValidation, "Email is required")
}
```

### 4. Log Strategically

```go
logger := gc.GetLogger()

// Info: Successful operations
logger.Info(ctx, "user created", "user_id", user.ID)

// Warn: Validation issues
logger.Warn(ctx, "invalid email", "email", email)

// Error: System issues
logger.Error(ctx, "database error", err)
```

### 5. Use Path Parameters for IDs

```go
// Good: Type-safe path parameter extraction
gc.RegisterHandler(server, "GET /api/users/{id}", getUserHandler)

// Avoid: Query parameters for resource IDs
// GET /api/users?id=123
```

## Contributing

Contributions are welcome! Please ensure:

- ✅ All tests pass: `go test ./...`
- ✅ No linting issues: `golangci-lint run`
- ✅ Security scanning passes: `gosec ./...`
- ✅ Code is well-documented

## License

MIT License - see [LICENSE.md](LICENSE.md) for details

## Resources

- **Framework Documentation**: [Package docs](https://pkg.go.dev/github.com/huijiro/go-charge)
- **Todo API Example**: [`examples/todo-api/README.md`](examples/todo-api/README.md)
- **Hello World Example**: [`examples/hello-world`](examples/hello-world)
- **CI/CD Pipeline**: [`.github/CI-CD.md`](.github/CI-CD.md)

## Roadmap

Future enhancements (non-breaking):

- [ ] Request body size limits
- [ ] CORS middleware
- [ ] Rate limiting middleware
- [ ] OpenAPI/Swagger generation
- [ ] Graceful shutdown helpers
- [ ] Request timeout utilities
- [ ] Form data support

## FAQ

**Q: Why no external dependencies?**  
A: Standard library is production-ready and stable. Fewer dependencies = fewer security concerns and faster compile times.

**Q: Does GoCharge support middleware?**  
A: Yes! GoCharge provides a composable middleware chain with built-in logging and recovery middleware.

**Q: Can I customize error responses?**  
A: Absolutely. Implement the `ErrorEncoder` interface and pass it to `SetErrorEncoder()`.

**Q: How do I handle database transactions?**  
A: Store transaction context in request context using `WithValue()` and retrieve in handlers.

**Q: Does GoCharge support WebSockets?**  
A: Not natively, but you can use `http.Upgrader` within handlers since handlers have access to `http.ResponseWriter` and `http.Request`.

## Support

- 🐛 [Report Issues](https://github.com/huijiro/go-charge/issues)
- 💡 [Request Features](https://github.com/huijiro/go-charge/discussions)
- 📚 [Read Documentation](examples/todo-api/README.md)

---

**Built with ❤️ for developers who value type safety and simplicity.**
