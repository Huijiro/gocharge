package gocharge_test

import (
	"encoding/json"
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
