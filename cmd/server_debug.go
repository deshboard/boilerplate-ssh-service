package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/goph/stdlib/expvar"
	"github.com/goph/stdlib/net/http/pprof"
	"github.com/goph/stdlib/x/net/trace"
)

// newDebugServer creates a new debug and health check server.
func newDebugServer(a *application) serverz.Server {
	handler := http.NewServeMux()

	// Add health checks
	handler.Handle("/healthz", a.healthCollector.Handler(healthz.LivenessCheck))
	handler.Handle("/readiness", a.healthCollector.Handler(healthz.ReadinessCheck))

	if a.config.Debug {
		// This is probably okay, as this service should not be exposed to public in the first place.
		trace.SetAuth(trace.NoAuth)

		expvar.RegisterRoutes(handler)
		pprof.RegisterRoutes(handler)
		trace.RegisterRoutes(handler)
	}

	// Register application specific debug routes (like metrics, etc)
	registerDebugRoutes(a, handler)

	return &serverz.AppServer{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(a.logger)), "health: ", 0),
		},
		Name:   "debug",
		Addr:   serverz.NewAddr("tcp", a.config.DebugAddr),
		Logger: a.logger,
	}
}
