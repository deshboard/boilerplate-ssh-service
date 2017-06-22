package main

import (
	"io"
	"log"
	"net/http"

	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newHealthServer(logWriter io.Writer, checkerCollector healthz.Collector) serverz.Server {
	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", checkerCollector.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", checkerCollector.Handler(healthz.ReadinessCheck))
	healthHandler.Handle("/metrics", promhttp.Handler())

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: log.New(logWriter, "health: ", 0),
		},
		Name: "health",
	}
}
