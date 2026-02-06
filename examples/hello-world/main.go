package main

import (
	"context"

	gc "github.com/huijiro/go-charge"
)

// Response types
type GreetingRequest struct {
	Name string `json:"name"`
}

type GreetingResponse struct {
	Message string `json:"message"`
}

// Handler
func greetingHandler(ctx context.Context, w gc.Response[GreetingResponse], r gc.Request[GreetingRequest]) error {
	req, err := r.JSON()
	if err != nil {
		return gc.NewAppError(gc.ErrBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return gc.NewAppError(gc.ErrValidation, "Name is required")
	}

	logger := gc.GetLogger()
	logger.Info(ctx, "greeting request", "name", req.Name)

	_, err = w.JSON(GreetingResponse{
		Message: "Hello, " + req.Name + "!",
	})
	return err
}

func main() {
	// Create server
	server := gc.New(":8080")
	logger := gc.GetLogger()

	// Register handler
	gc.RegisterHandler(server, "/greet", greetingHandler)

	// Start server
	logger.Info(context.Background(), "server starting on :8080")
	server.ListenAndServe()
}
