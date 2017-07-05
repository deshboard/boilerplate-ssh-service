package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/goph/serverz/named"
)

// newHealthServer creates a new health server and a status checker.
//
// The status checher can be used to manually mark the service unhealthy.
func newHealthServer(appCtx *application) (serverz.Server, *healthz.StatusChecker) {
	status := healthz.NewStatusChecker(healthz.Healthy)
	appCtx.healthCollector.RegisterChecker(healthz.ReadinessCheck, status)

	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", appCtx.healthCollector.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", appCtx.healthCollector.Handler(healthz.ReadinessCheck))

	if mReporter, ok := appCtx.metricsReporter.(interface {
		// HTTPHandler provides a scrape handler.
		HTTPHandler() http.Handler
	}); ok {
		healthHandler.Handle("/metrics", mReporter.HTTPHandler())
	}

	return &named.Server{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(appCtx.logger)), "health: ", 0),
		},
		ServerName: "health",
	}, status
}
