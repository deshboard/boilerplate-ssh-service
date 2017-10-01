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
func newDebugServer(app *application) serverz.Server {
	handler := http.NewServeMux()

	// Add health checks
	handler.Handle("/healthz", app.healthCollector.Handler(healthz.LivenessCheck))
	handler.Handle("/readiness", app.healthCollector.Handler(healthz.ReadinessCheck))

	if app.config.Debug {
		// This is probably okay, as this service should not be exposed to public in the first place.
		trace.SetAuth(trace.NoAuth)

		expvar.RegisterRoutes(handler)
		pprof.RegisterRoutes(handler)
		trace.RegisterRoutes(handler)
	}

	// Register application specific debug routes (like metrics, etc)
	registerDebugRoutes(app, handler)

	return &serverz.AppServer{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(log.With(app.Logger(), "server", "debug"))), "", 0),
		},
		Name:   "debug",
		Addr:   serverz.NewAddr("tcp", app.config.DebugAddr),
		Logger: app.Logger(),
	}
}
