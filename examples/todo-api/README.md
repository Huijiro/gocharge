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
**GET /health**

Health check endpoint for monitoring.

Response:
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

### Create Todo
**POST /api/todos**

Create a new todo item.

Request:
```json
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

Error response (400 Bad Request):
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Title is required"
}
```

### List Todos
**GET /api/todos**

List all todo items.

Response:
```json
{
  "todos": [
    {
      "id": "todo-1",
      "title": "Buy milk",
      "done": false
    },
    {
      "id": "todo-2",
      "title": "Learn GoCharge",
      "done": true
    }
  ],
  "count": 2
}
```

### Get Todo by ID
**GET /api/todos/{id}**

Get a single todo item by ID. Uses path parameter extraction.

Path Parameters:
- `id` (string) - The todo ID (e.g., "todo-1")

Response (200 OK):
```json
{
  "id": "todo-1",
  "title": "Buy milk",
  "done": false
}
```

Error response (404 Not Found):
```json
{
  "code": "NOT_FOUND",
  "message": "Todo not found"
}
```

### Update Todo
**PUT /api/todos/{id}**

Update an existing todo item. Uses path parameter extraction.

Path Parameters:
- `id` (string) - The todo ID (e.g., "todo-1")

Request:
```json
{
  "title": "Buy milk and eggs",
  "done": true
}
```

Response (200 OK):
```json
{
  "id": "todo-1",
  "title": "Buy milk and eggs",
  "done": true
}
```

Error response (404 Not Found):
```json
{
  "code": "NOT_FOUND",
  "message": "Todo not found"
}
```

Error response (400 Bad Request):
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Title is required"
}
```

### Delete Todo
**DELETE /api/todos/{id}**

Delete a todo item. Uses path parameter extraction.

Path Parameters:
- `id` (string) - The todo ID (e.g., "todo-1")

Response (204 No Content) - Empty body on success

Error response (404 Not Found):
```json
{
  "code": "NOT_FOUND",
  "message": "Todo not found"
}
```

## Example Usage

### Health check
```bash
curl http://localhost:8080/health
```

### Create a todo
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn GoCharge"}'
```

Response:
```json
{
  "id": "todo-1",
  "title": "Learn GoCharge",
  "done": false
}
```

### List all todos
```bash
curl http://localhost:8080/api/todos
```

Response:
```json
{
  "todos": [
    {
      "id": "todo-1",
      "title": "Learn GoCharge",
      "done": false
    }
  ],
  "count": 1
}
```

### Get a single todo by ID
```bash
curl http://localhost:8080/api/todos/todo-1
```

Response:
```json
{
  "id": "todo-1",
  "title": "Learn GoCharge",
  "done": false
}
```

### Update a todo
```bash
curl -X PUT http://localhost:8080/api/todos/todo-1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn GoCharge Framework", "done": true}'
```

Response:
```json
{
  "id": "todo-1",
  "title": "Learn GoCharge Framework",
  "done": true
}
```

### Delete a todo
```bash
curl -X DELETE http://localhost:8080/api/todos/todo-1
```

Response: 204 No Content (empty body)

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

### Test 404 Not Found
```bash
curl http://localhost:8080/api/todos/nonexistent
```

Error response (404 Not Found):
```json
{
  "code": "NOT_FOUND",
  "message": "Todo not found"
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
- `GetTodoByIDHandler` - Get single todo by ID (uses path parameter extraction)
- `UpdateTodoHandler` - Update a todo (uses path parameter extraction)
- `DeleteTodoHandler` - Delete a todo (uses path parameter extraction)

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

### 7. Path Parameter Extraction
```go
// Extract string parameter from path
id, err := r.PathParam().String("id")
if err != nil {
    return gc.NewAppError(gc.ErrBadRequest, err.Error())
}

// Extract and validate UUID parameter
uuid, err := r.PathParam().UUID("id")

// Extract integer parameter
count, err := r.PathParam().Int("count")

// Extract int64 parameter
itemID, err := r.PathParam().Int64("itemId")
```

### 8. Method-Aware Route Registration
```go
// Register handlers with HTTP method prefix (Go 1.22+)
gc.RegisterHandler(server, "GET /api/todos", listHandler)
gc.RegisterHandler(server, "POST /api/todos", createHandler)
gc.RegisterHandler(server, "GET /api/todos/{id}", getHandler)
gc.RegisterHandler(server, "PUT /api/todos/{id}", updateHandler)
gc.RegisterHandler(server, "DELETE /api/todos/{id}", deleteHandler)
```

## Next Steps

To extend this example:

1. **Database**: Replace in-memory storage with database/sql
2. **Authentication**: Add JWT middleware
3. **Validation**: Add more comprehensive validation with field-level rules
4. **Testing**: Add more comprehensive integration tests
5. **Custom Errors**: Implement custom error encoder with detailed error fields
6. **Rate Limiting**: Add rate limiting middleware
7. **CORS**: Add CORS middleware for frontend integration
8. **Pagination**: Add pagination support to list endpoint
9. **Filtering**: Add filtering and sorting to list endpoint
10. **Batch Operations**: Add batch create/delete endpoints

## Learning Resources

- [GoCharge Documentation](../..)
- [Sub-routes Guide](../../.github/CI-CD.md)
- [CI/CD Pipeline](../../.github/CI-CD.md)
