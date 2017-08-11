package main

import (
	stdlog "log"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz/aio"
	"github.com/goph/stdlib/net"
)

// newHTTPServer creates the main server instance for the service.
func newHTTPServer(appCtx *application) *aio.Server {
	serviceChecker := healthz.NewTCPChecker(appCtx.config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	appCtx.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	return &aio.Server{
		Server: &http.Server{
			Handler:  mux,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(appCtx.logger)), "http: ", 0),
		},
		Name: "http",
		Addr: net.ResolveVirtualAddr("tcp", appCtx.config.ServiceAddr),
	}
}
