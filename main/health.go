package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// newHealthServer creates a new health server and a status checker
func newHealthServer(logger log.Logger, healthCollector healthz.Collector) (serverz.Server, *healthz.StatusChecker) {
	status := healthz.NewStatusChecker(healthz.Healthy)
	healthCollector.RegisterChecker(healthz.ReadinessCheck, status)

	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", healthCollector.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", healthCollector.Handler(healthz.ReadinessCheck))
	healthHandler.Handle("/metrics", promhttp.Handler())

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "health: ", 0),
		},
		Name: "health",
	}, status
}
