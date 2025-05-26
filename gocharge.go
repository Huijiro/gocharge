package gocharge

import (
	"net/http"
)

type Server struct {
	http.Server
	Handler *http.ServeMux
}

func New(addr string) *Server {
	mux := http.NewServeMux()
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: mux,
		},
		Handler: mux,
	}
}
