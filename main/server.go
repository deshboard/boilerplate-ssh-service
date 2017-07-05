package main

import (
	stdlog "log"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/goph/serverz/aio"
	opentracing "github.com/opentracing/opentracing-go"
)

// newServer creates the main server instance for the service.
func newServer(config *configuration, logger log.Logger, errorHandler emperror.Handler, tracer opentracing.Tracer, healthCollector healthz.Collector, metricsReporter interface{}) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	return &aio.Server{
		Server: &http.Server{
			Handler:  mux,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "http: ", 0),
		},
		ServerName: "http",
	}
}
