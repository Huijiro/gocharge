package gocharge

import (
	"net/http"
)

type Server struct {
	Internal *http.Server
}

type Options struct {
	Address string
}

func (s *Server) RegisterHandler(path string, handler HandlerFunc[any, any]) {
	s.getDefaultMux().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		resp := Response[any]{
			ResponseWriter: w,
		}
		req := Request[any]{
			Request: *r,
		}

		if err := handler(resp, req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	},
	)
}

func (s *Server) Get(path string, handler HandlerFunc[any, any]) {
	s.RegisterHandler("GET "+path, handler)
}

func (s *Server) Post(path string, handler HandlerFunc[any, any]) {
	s.RegisterHandler("POST "+path, handler)
}

func (s *Server) Put(path string, handler HandlerFunc[any, any]) {
	s.RegisterHandler("PUT "+path, handler)
}

func (s *Server) Patch(path string, handler HandlerFunc[any, any]) {
	s.RegisterHandler("PATCH "+path, handler)
}

func (s *Server) Delete(path string, handler HandlerFunc[any, any]) {
	s.RegisterHandler("DELETE "+path, handler)
}

func (s *Server) ListenAndServe() error {
	return s.Internal.ListenAndServe()
}

func (s *Server) getDefaultMux() *http.ServeMux {
	return s.Internal.Handler.(*http.ServeMux)
}

func (s *Server) Run() error {
	return s.Internal.ListenAndServe()
}

func New(options *Options) *Server {
	return &Server{
		Internal: &http.Server{
			Addr:    options.Address,
			Handler: http.NewServeMux(),
		},
	}
}
