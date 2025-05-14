package gocharge_test

import (
	"log"
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
		&gocharge.Options{
			Address: addr,
		},
	)

	go server.ListenAndServe()

	// If the server ever takes more than 300 ms to start my code is shit and I should rethink it.
	time.Sleep(300 * time.Millisecond)

	m.Run()
}

func TestServerRouteRegister(t *testing.T) {
	server.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	resp, err := http.Get("http://" + addr + "/test")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("status code is not 200")
	}
}

func TestServerHealth(t *testing.T) {
	resp, err := http.Get("http://" + addr + "/_health")

	log.Printf("Checking server health on: %v", addr)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("status code is not 200")
	}
}
