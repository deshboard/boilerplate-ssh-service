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
	_trace "github.com/goph/stdlib/x/net/trace"
	"golang.org/x/net/trace"
)

// newDebugServer creates a new debug and health check server.
func newDebugServer(appCtx *application) serverz.Server {
	handler := http.NewServeMux()

	// Add health checks
	handler.Handle("/healthz", appCtx.healthCollector.Handler(healthz.LivenessCheck))
	handler.Handle("/readiness", appCtx.healthCollector.Handler(healthz.ReadinessCheck))

	if appCtx.config.Debug {
		trace.AuthRequest = _trace.NoAuth

		expvar.RegisterRoutes(handler)
		pprof.RegisterRoutes(handler)
		_trace.RegisterRoutes(handler)
	}

	// Register application specific debug routes (like metrics, etc)
	registerDebugRoutes(appCtx, handler)

	return &serverz.AppServer{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(appCtx.logger)), "health: ", 0),
		},
		Name:   "debug",
		Addr:   serverz.NewAddr("tcp", appCtx.config.DebugAddr),
		Logger: appCtx.logger,
	}
}
