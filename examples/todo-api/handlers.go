package main

import (
	"context"

	gc "github.com/huijiro/go-charge"
)

// ============= Handler Helpers =============

// These are template handlers that show patterns for
// implementing GET, PUT, DELETE operations.
//
// In a real application with URL path parameters,
// you would extract the todo ID from the URL and
// perform the corresponding operation.

// GetTodoByIDHandler - Get a single todo by ID
//
// Pattern:
//
//	GET /api/todos/{id}
//
// In practice, you would:
//  1. Extract the ID from the URL path
//  2. Look up the todo in the store
//  3. Return 404 if not found
//  4. Return the todo if found
func GetTodoByIDHandler(store *TodoStore) gc.HandlerFunc[TodoResponse, string] {
	return func(ctx context.Context, w gc.Response[TodoResponse], r gc.Request[string]) error {
		logger := gc.GetLogger()

		// TODO: Extract ID from URL path
		// For now, just demonstrate the error pattern
		logger.Warn(ctx, "get todo handler - implement URL extraction")

		return gc.NewAppError(gc.ErrNotFound, "Todo not found")
	}
}

// UpdateTodoHandler - Update a todo
//
// Pattern:
//
//	PUT /api/todos/{id}
//
// In practice, you would:
//  1. Extract the ID from the URL path
//  2. Parse the request body
//  3. Validate the input
//  4. Update the todo in the store
//  5. Return the updated todo
func UpdateTodoHandler(store *TodoStore) gc.HandlerFunc[TodoResponse, UpdateTodoRequest] {
	return func(ctx context.Context, w gc.Response[TodoResponse], r gc.Request[UpdateTodoRequest]) error {
		logger := gc.GetLogger()

		req, err := r.JSON()
		if err != nil {
			logger.Warn(ctx, "invalid json in update todo request")
			return gc.NewAppError(gc.ErrBadRequest, "Invalid request body")
		}

		// Validate
		if req.Title == "" {
			logger.Warn(ctx, "empty title in update todo request")
			return gc.NewAppError(gc.ErrValidation, "Title is required")
		}

		// TODO: Extract ID from URL path
		// For now, just demonstrate the error pattern
		logger.Warn(ctx, "update todo handler - implement URL extraction")

		return gc.NewAppError(gc.ErrNotFound, "Todo not found")
	}
}

// DeleteTodoHandler - Delete a todo
//
// Pattern:
//
//	DELETE /api/todos/{id}
//
// In practice, you would:
//  1. Extract the ID from the URL path
//  2. Delete the todo from the store
//  3. Return 204 No Content on success
//  4. Return 404 if not found
func DeleteTodoHandler(store *TodoStore) gc.HandlerFunc[MessageResponse, string] {
	return func(ctx context.Context, w gc.Response[MessageResponse], r gc.Request[string]) error {
		logger := gc.GetLogger()

		// TODO: Extract ID from URL path
		// For now, just demonstrate the error pattern
		logger.Warn(ctx, "delete todo handler - implement URL extraction")

		return gc.NewAppError(gc.ErrNotFound, "Todo not found")
	}
}
