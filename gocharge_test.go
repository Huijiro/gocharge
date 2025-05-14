package gocharge_test

import (
	"errors"
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
	server.Get("/", func(w gocharge.Response[any], r gocharge.Request[any]) error {
		return nil
	})

	resp, err := http.Get("http://" + addr + "/")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code is not 200, server responded with: %s", resp.Status)
	}
}

func TestError(t *testing.T) {
	server.Get("/should-error", func(w gocharge.Response[any], r gocharge.Request[any]) error {
		return errors.New("test")
	})

	resp, err := http.Get("http://" + addr + "/should-error")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code is not 500, server responded with: %s", resp.Status)
	}
}
