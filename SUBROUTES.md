# Sub-Routes Guide

GoCharge supports any sub-route pattern through the standard `RegisterHandler` function. Since it uses Go's standard `http.ServeMux` under the hood, all paths are supported.

## Basic Patterns

### Simple Sub-routes

```go
// Base routes
gocharge.RegisterHandler(server, "/api/users", listUsersHandler)
gocharge.RegisterHandler(server, "/api/posts", listPostsHandler)

// Nested routes
gocharge.RegisterHandler(server, "/api/users/{id}", getUserHandler)
gocharge.RegisterHandler(server, "/api/posts/{id}", getPostHandler)
```

### Version Prefixes

```go
// Version 1 API
gocharge.RegisterHandler(server, "/api/v1/users", v1ListUsers)
gocharge.RegisterHandler(server, "/api/v1/posts", v1ListPosts)

// Version 2 API (with improvements)
gocharge.RegisterHandler(server, "/api/v2/users", v2ListUsers)
gocharge.RegisterHandler(server, "/api/v2/posts", v2ListPosts)
```

## Hierarchical Resources

### Two-level Nesting

```go
// User's posts endpoint
gocharge.RegisterHandler(server, "/api/users/{userId}/posts", getUserPosts)

// User's comments endpoint
gocharge.RegisterHandler(server, "/api/users/{userId}/comments", getUserComments)
```

### Three or More Levels

```go
// Organization → Project → Tasks
gocharge.RegisterHandler(
    server,
    "/api/organizations/{orgId}/projects/{projectId}/tasks",
    listOrgProjectTasks,
)

// Organization → Team → Members
gocharge.RegisterHandler(
    server,
    "/api/organizations/{orgId}/teams/{teamId}/members",
    listTeamMembers,
)
```

## Domain-Organized Routes

Group routes by functional domain:

```go
// User domain
gocharge.RegisterHandler(server, "/api/users", listUsers)
gocharge.RegisterHandler(server, "/api/users/{id}", getUser)

// Profile domain
gocharge.RegisterHandler(server, "/api/profile", getProfile)
gocharge.RegisterHandler(server, "/api/profile/settings", getSettings)

// Admin domain
gocharge.RegisterHandler(server, "/admin/users", adminListUsers)
gocharge.RegisterHandler(server, "/admin/settings", adminSettings)
```

## RESTful API Pattern

Typical REST structure with sub-routes:

```go
// List all items
gocharge.RegisterHandler(server, "/api/items", func(ctx context.Context, w gocharge.Response[ItemListResponse], r gocharge.Request[string]) error {
    // GET /api/items
})

// Create new item
gocharge.RegisterHandler(server, "/api/items", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[ItemCreateRequest]) error {
    // POST /api/items
})

// Get single item
gocharge.RegisterHandler(server, "/api/items/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[string]) error {
    // GET /api/items/{id}
})

// Update item
gocharge.RegisterHandler(server, "/api/items/{id}", func(ctx context.Context, w gocharge.Response[ItemResponse], r gocharge.Request[ItemUpdateRequest]) error {
    // PUT /api/items/{id}
})

// Delete item
gocharge.RegisterHandler(server, "/api/items/{id}", func(ctx context.Context, w gocharge.Response[struct{}], r gocharge.Request[string]) error {
    // DELETE /api/items/{id}
})
```

## Error Handling in Sub-routes

Sub-routes automatically handle errors consistently:

```go
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

gocharge.RegisterHandler(server, "/api/v2/users/{id}", func(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[GetUserRequest]) error {
    req, err := r.JSON()
    if err != nil {
        return gocharge.NewAppError(gocharge.ErrBadRequest, "Invalid request body")
    }

    if req.ID == "" {
        return gocharge.NewAppError(gocharge.ErrValidation, "User ID is required")
    }

    user, found := database.GetUser(req.ID)
    if !found {
        return gocharge.NewAppError(gocharge.ErrNotFound, "User not found")
    }

    // Success
    _, err = w.Status(gocharge.StatusOK).JSON(UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    })
    return err
})
```

All errors are automatically encoded with proper HTTP status codes and JSON responses.

## Path Variations

Any path pattern works:

```go
// Webhooks
gocharge.RegisterHandler(server, "/webhooks/github", githubWebhookHandler)
gocharge.RegisterHandler(server, "/webhooks/stripe", stripeWebhookHandler)

// Admin
gocharge.RegisterHandler(server, "/admin/dashboard", adminDashboard)
gocharge.RegisterHandler(server, "/admin/users", adminUsers)

// GraphQL
gocharge.RegisterHandler(server, "/graphql", graphqlHandler)
gocharge.RegisterHandler(server, "/graphql/playground", playgroundHandler)

// Internal
gocharge.RegisterHandler(server, "/internal/metrics", metricsHandler)
gocharge.RegisterHandler(server, "/internal/health", healthCheckHandler)

// Public
gocharge.RegisterHandler(server, "/public/assets", serveAssets)
gocharge.RegisterHandler(server, "/public/docs", serveDocs)
```

## Shared Middleware for Routes

Group routes that need the same middleware:

```go
// Logging middleware
loggingMW := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger := gocharge.GetLogger()
        logger.Info(r.Context(), "incoming_request", "path", r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

// Auth middleware
authMW := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        // Verify token...
        next.ServeHTTP(w, r)
    })
}

// Public routes with logging only
publicChain := gocharge.NewChain()
publicChain.Use(loggingMW)

// Protected routes with auth + logging
protectedChain := gocharge.NewChain()
protectedChain.Use(loggingMW)
protectedChain.Use(authMW)

// When you need per-route middleware handling, handlers can check context
gocharge.RegisterHandler(server, "/api/public/data", publicDataHandler)
gocharge.RegisterHandler(server, "/api/protected/data", protectedDataHandler)
```

## Organizing Handler Code

For larger applications, organize handlers by sub-route:

```go
// handlers/users.go
package handlers

type UsersHandler struct {
    db *sql.DB
}

func (h *UsersHandler) List(ctx context.Context, w gocharge.Response[ListResponse], r gocharge.Request[string]) error {
    // List users
}

func (h *UsersHandler) Get(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[GetRequest]) error {
    // Get single user
}

func (h *UsersHandler) Create(ctx context.Context, w gocharge.Response[UserResponse], r gocharge.Request[CreateRequest]) error {
    // Create user
}

// handlers/posts.go
package handlers

type PostsHandler struct {
    db *sql.DB
}

func (h *PostsHandler) List(ctx context.Context, w gocharge.Response[ListResponse], r gocharge.Request[string]) error {
    // List posts
}

// main.go
package main

func main() {
    server := gocharge.New(":8080")
    
    usersH := &handlers.UsersHandler{db: db}
    postsH := &handlers.PostsHandler{db: db}
    
    // Register user routes
    gocharge.RegisterHandler(server, "/api/users", usersH.List)
    gocharge.RegisterHandler(server, "/api/users/{id}", usersH.Get)
    
    // Register post routes
    gocharge.RegisterHandler(server, "/api/posts", postsH.List)
    gocharge.RegisterHandler(server, "/api/posts/{id}", postsH.Get)
    
    server.ListenAndServe()
}
```

## Best Practices

1. **Be Consistent**: Use consistent naming patterns across similar routes
   - `/api/v1/users` and `/api/v1/posts` (consistent structure)
   - Not: `/api/v1/users` and `/api/items` (inconsistent)

2. **Use Meaningful Paths**: Make the path structure reflect your API semantics
   - `/api/organizations/{orgId}/projects/{projectId}/tasks` (clear hierarchy)
   - Not: `/api/o/{o}/p/{p}/t` (unclear)

3. **Version Early**: Plan for versioning from the start
   - `/api/v1/...` (versioned from start)
   - Not: `/api/...` then try to add versioning later

4. **Error Handling**: Always return typed errors from handlers
   ```go
   if err != nil {
       return gocharge.NewAppError(gocharge.ErrNotFound, "Resource not found")
   }
   ```

5. **Keep Handlers Simple**: Move business logic to separate packages
   ```go
   // Good: handler calls service
   user, err := userService.Get(req.ID)
   
   // Bad: handler contains all logic
   rows, err := db.Query(...) // Don't do this in handler
   ```

## Examples

See `subroute_test.go` for comprehensive examples of:
- Basic sub-routes
- Version prefixes
- Hierarchical resources
- Multiple nesting levels
- Error handling
- RESTful patterns
- Path variations
