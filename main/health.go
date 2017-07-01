package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
)

// newHealthServer creates a new health server and a status checker.
//
// The status checher can be used to manually mark the service unhealthy.
func newHealthServer(logger log.Logger, healthCollector healthz.Collector, metricsReporter interface{}) (serverz.Server, *healthz.StatusChecker) {
	status := healthz.NewStatusChecker(healthz.Healthy)
	healthCollector.RegisterChecker(healthz.ReadinessCheck, status)

	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", healthCollector.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", healthCollector.Handler(healthz.ReadinessCheck))

	if mReporter, ok := metricsReporter.(interface {
		// HTTPHandler provides a scrape handler.
		HTTPHandler() http.Handler
	}); ok {
		healthHandler.Handle("/metrics", mReporter.HTTPHandler())
	}

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "health: ", 0),
		},
		Name: "health",
	}, status
}
