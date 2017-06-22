package main

import (
	_ "expvar"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/goph/serverz"

	"golang.org/x/net/trace"
)

func init() {
	// This is probably OK as the service runs in Docker
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
}

func newDebugServer(logWriter io.Writer) serverz.Server {
	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  http.DefaultServeMux,
			ErrorLog: log.New(logWriter, "debug: ", 0),
		},
		Name: "debug",
	}
}
