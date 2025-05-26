package gocharge_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

type TypeResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

func TestHandler(t *testing.T) {
	testResponse := &TypeResponse{
		Message: "Hello World",
		Data:    "This is a test",
	}

	gocharge.RegisterHandler(server, "/testHandler", func(w gocharge.Response[TypeResponse], r gocharge.Request[string]) error {
		w.JSON(*testResponse)
		return nil
	})

	response, err := http.Get("http://localhost:8080/testHandler")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	serverResponse := &TypeResponse{}
	err = json.NewDecoder(response.Body).Decode(&serverResponse)

	if !reflect.DeepEqual(serverResponse, testResponse) {
		t.Errorf("Expected: %v, got: %v", testResponse, serverResponse)
	}
}

func TestStringHandler(t *testing.T) {
	testResponse := "Hello World"

	gocharge.RegisterHandler(server, "/testStringHandler", func(w gocharge.Response[string], r gocharge.Request[string]) error {
		w.Write([]byte(testResponse))
		return nil
	})

	response, err := http.Get("http://localhost:8080/testStringHandler")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if string(body) != testResponse {
		t.Errorf("Expected: %v, got: %v", testResponse, string(body))
	}
}

func TestRequest(t *testing.T) {
	testResponse := &TypeResponse{
		Message: "Hello World",
		Data:    "This is a test",
	}

	gocharge.RegisterHandler(server, "/testRequest", func(w gocharge.Response[TypeResponse], r gocharge.Request[TypeResponse]) error {
		req, err := r.JSON()
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		w.JSON(*req)

		return nil
	})

	body, err := json.Marshal(testResponse)

	request, err := http.NewRequest("GET", "http://localhost:8080/testRequest", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	serverResponse := &TypeResponse{}
	err = json.NewDecoder(response.Body).Decode(&serverResponse)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if !reflect.DeepEqual(serverResponse, testResponse) {
		t.Errorf("Expected: %v, got: %v", testResponse, serverResponse)
	}
}
