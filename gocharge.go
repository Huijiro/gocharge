package gocharge

import (
	"net/http"
)

// Server is the main HTTP server type for gocharge
// It extends http.Server with additional typed features
type Server struct {
	http.Server
	Handler         *http.ServeMux
	ErrorEncoder    ErrorEncoder
	Logger          Logger
	MiddlewareChain *Chain
}

// New creates a new Server instance with default settings
func New(addr string) *Server {
	mux := http.NewServeMux()
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: mux,
		},
		Handler:         mux,
		ErrorEncoder:    &DefaultErrorEncoder{},
		Logger:          GetLogger(),
		MiddlewareChain: NewChain(),
	}
}

// SetErrorEncoder sets the error encoder for the server
// This determines how errors are converted to JSON responses
func (s *Server) SetErrorEncoder(enc ErrorEncoder) {
	s.ErrorEncoder = enc
}

// SetLogger sets the logger for the server
func (s *Server) SetLogger(l Logger) {
	s.Logger = l
}

// Middleware returns the middleware chain for this server
// Users can call Use() on the returned chain to add middleware
func (s *Server) Middleware() *Chain {
	return s.MiddlewareChain
}
