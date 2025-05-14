package gocharge

import "net/http"

type Server struct {
	DefaultMux *http.ServeMux
	Internal   *http.Server
}

type Options struct {
	Address string
}

func (s *Server) Get(path string, handler http.HandlerFunc) {
	s.DefaultMux.HandleFunc("GET "+path, handler)
}

func (s *Server) Post(path string, handler http.HandlerFunc) {
	s.DefaultMux.HandleFunc("POST "+path, handler)
}

func (s *Server) Put(path string, handler http.HandlerFunc) {
	s.DefaultMux.HandleFunc("PUT "+path, handler)
}

func (s *Server) Delete(path string, handler http.HandlerFunc) {
	s.DefaultMux.HandleFunc("DELETE "+path, handler)
}

func (s *Server) ListenAndServe() error {
	return s.Internal.ListenAndServe()
}

func New(options *Options) *Server {
	mux := http.NewServeMux()
	return &Server{
		DefaultMux: mux,
		Internal: &http.Server{
			Addr:    options.Address,
			Handler: mux,
		},
	}
}
