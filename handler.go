package gocharge

import (
	"log"
	"net/http"
)

type WTypes interface {
	string | any
}

type RTypes interface {
	string | any
}

type HandlerFunc[W WTypes, R RTypes] func(w Response[W], r Request[R]) error

func RegisterHandler[W WTypes, R RTypes](s *Server, path string, handler HandlerFunc[W, R]) {
	log.Printf("Registering handler: %s", path)
	s.Handler.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		response := Response[W]{
			ResponseWriter: w,
			Data:           *new(W),
			StatusCode:     http.StatusOK,
		}

		request := Request[R]{
			Request: *r,
			Data:    *new(R),
		}

		handler(response, request)
	})
}
