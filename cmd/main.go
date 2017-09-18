package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/goph/stdlib/ext"
)

func main() {
	app, err := newApplication(
		configProvider,
		loggerProvider,
		errorHandlerProvider,
	)
	// Close resources even when there is an error
	defer app.Close()

	if err != nil {
		panic(err)
	}

	// Register error handler to recover from panics
	defer emperror.HandleRecover(app.errorHandler)

	// Create a new health collector
	healthCollector := healthz.Collector{}

	// Create a new application tracer
	tracer := newTracer(app.config, app.logger)
	defer ext.Close(tracer)

	// Application context
	app.healthCollector = healthCollector
	app.tracer = tracer

	status := healthz.NewStatusChecker(healthz.Healthy)
	healthCollector.RegisterChecker(healthz.ReadinessCheck, status)

	level.Info(app.logger).Log(
		"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
		"version", Version,
		"commit_hash", CommitHash,
		"build_date", BuildDate,
		"environment", app.config.Environment,
	)

	serverQueue := newServerQueue(app)
	defer serverQueue.Close()

	errChan := serverQueue.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		status.SetStatus(healthz.Unhealthy)
		level.Debug(app.logger).Log("msg", "error received from error channel")
		emperror.HandleIfErr(app.errorHandler, err)
	case s := <-signalChan:
		level.Info(app.logger).Log("msg", fmt.Sprintf("captured %v", s))
		status.SetStatus(healthz.Unhealthy)

		level.Debug(app.logger).Log(
			"msg", "shutting down with timeout",
			"timeout", app.config.ShutdownTimeout,
		)

		ctx, cancel := context.WithTimeout(context.Background(), app.config.ShutdownTimeout)

		err := serverQueue.Shutdown(ctx)
		if err != nil {
			app.errorHandler.Handle(err)
		}

		// Cancel context if shutdown completed earlier
		cancel()
	}
}
