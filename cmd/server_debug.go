package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz/aio"
	"github.com/goph/stdlib/expvar"
	"github.com/goph/stdlib/net"
	"github.com/goph/stdlib/net/http/pprof"
	_trace "github.com/goph/stdlib/x/net/trace"
	"golang.org/x/net/trace"
)

// newDebugServer creates a new debug and health check server.
func newDebugServer(appCtx *application) *aio.Server {
	handler := http.NewServeMux()

	// Add health checks
	handler.Handle("/healthz", appCtx.healthCollector.Handler(healthz.LivenessCheck))
	handler.Handle("/readiness", appCtx.healthCollector.Handler(healthz.ReadinessCheck))

	// Check if a (Prometheus) HTTP handler is available
	if h, ok := appCtx.metrics.(http.Handler); ok {
		level.Debug(appCtx.logger).Log(
			"msg", "Exposing Prometheus metrics",
			"server", "health",
		)

		handler.Handle("/metrics", h)
	}

	if appCtx.config.Debug {
		trace.AuthRequest = _trace.NoAuth

		expvar.RegisterRoutes(handler)
		pprof.RegisterRoutes(handler)
		_trace.RegisterRoutes(handler)
	}

	return &aio.Server{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(appCtx.logger)), "health: ", 0),
		},
		Name: "debug",
		Addr: net.ResolveVirtualAddr("tcp", appCtx.config.DebugAddr),
	}
}
