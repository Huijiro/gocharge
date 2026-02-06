# Todo API Example

A simple but complete TODO API example demonstrating GoCharge features and best practices.

## Features Demonstrated

- ✅ Type-safe request/response handlers
- ✅ Error handling with typed errors
- ✅ Structured logging
- ✅ Middleware (logging, recovery)
- ✅ In-memory data storage with concurrency safety
- ✅ Request validation
- ✅ RESTful API patterns
- ✅ Response chaining

## Building

```bash
cd examples/todo-api
go build -o todo-api .
```

## Running

```bash
./todo-api
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
```bash
GET /health
```

Response:
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

### Create Todo
```bash
POST /api/todos
Content-Type: application/json

{
  "title": "Buy milk"
}
```

Response (201 Created):
```json
{
  "id": "todo-1",
  "title": "Buy milk",
  "done": false
}
```

### List Todos
```bash
GET /api/todos
```

Response:
```json
{
  "todos": [
    {
      "id": "todo-1",
      "title": "Buy milk",
      "done": false
    }
  ],
  "count": 1
}
```

## Example Usage

### Create a todo
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn GoCharge"}'
```

### List todos
```bash
curl http://localhost:8080/api/todos
```

### Health check
```bash
curl http://localhost:8080/health
```

### Test validation error
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'
```

Error response (400 Bad Request):
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Title is required"
}
```

## Code Structure

### Domain Types
- `Todo` - Core domain model
- `TodoStore` - Thread-safe in-memory storage

### Request/Response Types
- `CreateTodoRequest` - Create todo request
- `TodoResponse` - Single todo response
- `TodoListResponse` - List of todos response
- `MessageResponse` - Generic message response
- `HealthResponse` - Health check response

### Handlers
- `healthHandler` - Health check endpoint
- `createTodoHandler` - Create new todo
- `listTodosHandler` - List all todos
- `getTodoHandler` - Get single todo (template)
- `deleteTodoHandler` - Delete todo (template)

### Middleware
- `LoggingMiddleware` - Logs all requests and responses
- `RecoveryMiddleware` - Catches panics and returns 500

## Key GoCharge Features Used

### 1. Type-Safe Handlers
```go
func createTodoHandler(store *TodoStore) gc.HandlerFunc[TodoResponse, CreateTodoRequest] {
    return func(ctx context.Context, w gc.Response[TodoResponse], r gc.Request[CreateTodoRequest]) error {
        // Type-safe request/response
    }
}
```

### 2. Typed Errors
```go
return gc.NewAppError(gc.ErrValidation, "Title is required")
```

### 3. Request Validation
```go
req, err := r.JSON()
if err != nil {
    return gc.NewAppError(gc.ErrBadRequest, "Invalid request body")
}
```

### 4. Structured Logging
```go
logger := gc.GetLogger()
logger.Info(ctx, "todo created", "todo_id", todo.ID)
```

### 5. Response Chaining
```go
_, err = w.Status(gc.StatusCreated).JSON(response)
return err
```

### 6. Middleware
```go
chain := server.Middleware()
chain.Use(middleware.LoggingMiddleware)
chain.Use(middleware.RecoveryMiddleware)
```

## Next Steps

To extend this example:

1. **Database**: Replace in-memory storage with database/sql
2. **Authentication**: Add JWT middleware
3. **Validation**: Add more comprehensive validation
4. **Testing**: Add integration tests
5. **URL Parameters**: Extract todo ID from path parameters
6. **Custom Errors**: Implement custom error encoder
7. **Rate Limiting**: Add rate limiting middleware
8. **CORS**: Add CORS middleware for frontend integration

## Learning Resources

- [GoCharge Documentation](../..)
- [Sub-routes Guide](../../.github/CI-CD.md)
- [CI/CD Pipeline](../../.github/CI-CD.md)
