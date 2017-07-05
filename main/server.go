package main

import (
	stdlog "log"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/goph/serverz/aio"
)

// newServer creates the main server instance for the service.
func newServer(app *application) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(app.config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	app.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	return &aio.Server{
		Server: &http.Server{
			Handler:  mux,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(app.logger)), "http: ", 0),
		},
		ServerName: "http",
	}
}
