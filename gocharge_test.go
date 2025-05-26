package gocharge_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	gocharge "github.com/huijiro/go-charge"
)

var addr string
var server *gocharge.Server

func TestMain(m *testing.M) {
	addr = ":8080"
	server = gocharge.New(
		addr,
	)

	// Wait for the server to start before cont
	go server.ListenAndServe()

	fmt.Println("Server started")

	// If the server ever takes more than 300 ms to start my code is shit and I should rethink it.
	time.Sleep(300 * time.Millisecond)

	m.Run()
}

func TestServer(t *testing.T) {
	server.Handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	response, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(body) != "Hello World" {
		t.Errorf("Expected: %v, got: %v", "Hello World", string(body))
	}
}
