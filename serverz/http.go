package serverz

import (
	"context"
	"net/http"
)

// HTTPServer wraps an *http.Server
type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer creates a new HTTPServer
func NewHTTPServer(server *http.Server) *HTTPServer {
	return &HTTPServer{server}
}

// Serve starts listening
func (s *HTTPServer) Serve() error {
	return s.server.ListenAndServe()
}

// Shutdown initiates graceful shutdown and returns an error if the context is canceled
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Close forcefully stops the server
func (s *HTTPServer) Close() error {
	return s.server.Close()
}
