package gocharge

import (
	"encoding/json"
	"net/http"
)

type Request[T any] struct {
	Data T
	http.Request
}

func (r *Request[T]) JSON() (*T, error) {
	err := json.NewDecoder(r.Body).Decode(&r.Data)

	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}
