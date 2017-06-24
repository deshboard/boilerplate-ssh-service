package main

import (
	_ "expvar"
	stdlog "log"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/serverz"
	"golang.org/x/net/trace"
)

func init() {
	// This is probably OK as the service runs in Docker
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
}

// newDebugServer creates a debug server
func newDebugServer(logger log.Logger) serverz.Server {
	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  http.DefaultServeMux,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "debug: ", 0),
		},
		Name: "debug",
	}
}
