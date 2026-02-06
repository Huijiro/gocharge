# GoCharge Examples

This directory contains example applications demonstrating GoCharge features and best practices.

## Examples

### 1. Hello World ([hello-world](./hello-world))

The simplest GoCharge application - a greeting API with validation and error handling.

**Key Features:**
- Type-safe request/response handlers
- Request validation
- Typed error handling
- Structured logging

**Quick Start:**
```bash
cd hello-world
go build -o hello-world .
./hello-world
curl -X POST http://localhost:8080/greet \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```

### 2. Todo API ([todo-api](./todo-api))

A complete TODO API example with multiple endpoints and in-memory storage.

**Key Features:**
- Multiple typed handlers
- In-memory data storage with concurrency safety
- Middleware (logging, recovery)
- RESTful API patterns
- Response chaining
- Handler templates for GET, PUT, DELETE

**Quick Start:**
```bash
cd todo-api
go build -o todo-api .
./todo-api
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk"}'
```

## Running All Examples

```bash
# Hello World
cd hello-world && go run . &

# In another terminal - Todo API
cd todo-api && go run . &
```

## Common Patterns

### 1. Creating a Typed Handler

```go
type CreateRequest struct {
    Field string `json:"field"`
}

type Response struct {
    Result string `json:"result"`
}

func handler(ctx context.Context, w gc.Response[Response], r gc.Request[CreateRequest]) error {
    req, err := r.JSON()
    if err != nil {
        return gc.NewAppError(gc.ErrBadRequest, "Invalid JSON")
    }
    
    _, err = w.JSON(Response{Result: "success"})
    return err
}
```

### 2. Validation with Errors

```go
if req.Field == "" {
    return gc.NewAppError(gc.ErrValidation, "Field is required")
}

if len(req.Field) > 100 {
    return gc.NewAppError(gc.ErrValidation, "Field too long")
}
```

### 3. Logging in Handlers

```go
logger := gc.GetLogger()
logger.Info(ctx, "processing request", "user_id", id)
logger.Warn(ctx, "validation failed", "reason", "empty field")
logger.Error(ctx, "database error", err, "table", "users")
```

### 4. Response Chaining

```go
_, err := w.Status(gc.StatusCreated).JSON(response)
return err

// Or with error response
_, err := w.Status(gc.StatusBadRequest).Error("INVALID", "Invalid input")
return err
```

### 5. Middleware Setup

```go
chain := server.Middleware()
chain.Use(middleware.LoggingMiddleware)
chain.Use(middleware.RecoveryMiddleware)
```

### 6. Custom Logger

```go
opts := &slog.HandlerOptions{Level: slog.LevelDebug}
handler := slog.NewTextHandler(os.Stdout, opts)
logger := gc.NewDefaultLoggerWithHandler(handler)
gc.SetLogger(logger)
```

## Testing Examples

### Using curl

```bash
# Health check
curl http://localhost:8080/health

# Create resource
curl -X POST http://localhost:8080/api/resource \
  -H "Content-Type: application/json" \
  -d '{"field": "value"}'

# Test error handling
curl -X POST http://localhost:8080/api/resource \
  -H "Content-Type: application/json" \
  -d '{}'
```

### Using httpie (nicer output)

```bash
# Install: https://httpie.io
http :8080/health
http POST :8080/api/resource field=value
```

## Next Steps

To extend these examples:

1. **Add Database** - Replace in-memory storage with PostgreSQL/SQLite
2. **Add Authentication** - Implement JWT middleware
3. **Add Tests** - Write integration tests for handlers
4. **Add More Endpoints** - Implement full CRUD operations
5. **Add Custom Middleware** - Create domain-specific middleware
6. **Add Configuration** - Load config from environment/files
7. **Add Metrics** - Integrate with Prometheus
8. **Add Docker** - Create Dockerfile for containerization

## Learning Resources

- [GoCharge Main Documentation](../)
- [API Documentation](https://pkg.go.dev/github.com/huijiro/go-charge)
- [Sub-routes Guide](..)
- [CI/CD Pipeline Documentation](../.github/CI-CD.md)

## Example Structure

Each example includes:

- `main.go` - Application code
- `handlers.go` - Handler implementations (in todo-api)
- `go.mod` - Module definition
- `README.md` - Documentation

## Common Issues

### "package github.com/huijiro/go-charge not found"

The go.mod files use `replace` directives to point to the parent directory. Make sure you're running from the example directory:

```bash
cd examples/hello-world
go run .
```

### Port already in use

By default, examples listen on `:8080`. If that port is busy:

1. Kill the existing process: `lsof -ti:8080 | xargs kill`
2. Or modify the port in the example code

## Contributing

To add a new example:

1. Create a new directory under `examples/`
2. Add `main.go` with your example
3. Add `go.mod` with proper replace directive
4. Add `README.md` with documentation
5. Test that it builds: `go build -o app .`
6. Test that it runs and handles requests correctly
