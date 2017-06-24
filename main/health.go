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

func newHealthServer(logger log.Logger, checkerCollector healthz.Collector) serverz.Server {
	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", checkerCollector.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", checkerCollector.Handler(healthz.ReadinessCheck))
	healthHandler.Handle("/metrics", promhttp.Handler())

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "health: ", 0),
		},
		Name: "health",
	}
}
