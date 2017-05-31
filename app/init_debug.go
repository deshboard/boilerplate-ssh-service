package main

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"

	"golang.org/x/net/trace"
)

func init() {
	// This is probably OK as the service runs in Docker
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
}
