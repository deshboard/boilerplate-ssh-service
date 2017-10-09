package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/fw"
	"github.com/goph/serverz"
	"github.com/goph/stdlib/expvar"
	"github.com/goph/stdlib/net/http/pprof"
	"github.com/goph/stdlib/x/net/trace"
)

// newDebugServer creates a new debug and health check server.
func newDebugServer(app *fw.Application) serverz.Server {
	handler := app.MustGet("debug_handler").(*http.ServeMux)
	config := app.MustGet("config").(*configuration)

	if config.Debug {
		// This is probably okay, as this service should not be exposed to public in the first place.
		trace.SetAuth(trace.NoAuth)

		expvar.RegisterRoutes(handler)
		pprof.RegisterRoutes(handler)
		trace.RegisterRoutes(handler)
	}

	return &serverz.AppServer{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(log.With(app.Logger(), "server", "debug"))), "", 0),
		},
		Name:   "debug",
		Addr:   serverz.NewAddr("tcp", config.DebugAddr),
		Logger: app.Logger(),
	}
}
