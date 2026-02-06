# Hello World Example

The simplest GoCharge application - a greeting API.

## Building

```bash
cd examples/hello-world
go build -o hello-world .
```

## Running

```bash
./hello-world
```

## Testing

### Make a greeting request
```bash
curl -X POST http://localhost:8080/greet \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```

Response:
```json
{
  "message": "Hello, World!"
}
```

### Test validation error
```bash
curl -X POST http://localhost:8080/greet \
  -H "Content-Type: application/json" \
  -d '{"name": ""}'
```

Response (400 Bad Request):
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Name is required"
}
```

### Test invalid JSON
```bash
curl -X POST http://localhost:8080/greet \
  -H "Content-Type: application/json" \
  -d 'not json'
```

Response (400 Bad Request):
```json
{
  "code": "BAD_REQUEST",
  "message": "Invalid request body"
}
```

## Key Concepts

This example demonstrates:

1. **Type-safe handlers** - Request and response types are generic
2. **Typed errors** - Using `gc.NewAppError` with error codes
3. **Request validation** - Checking for required fields
4. **Logging** - Using structured logging
5. **Response chaining** - Using Status/JSON/Error methods
