package main

import (
	stdlog "log"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	opentracing "github.com/opentracing/opentracing-go"
)

// newServer creates the main server instance for the service
func newServer(config *Configuration, logger log.Logger, checkerCollector healthz.Collector) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	checkerCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)
	_ = opentracing.GlobalTracer()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  mux,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "http: ", 0),
		},
		Name: "http",
	}
}
