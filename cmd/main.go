package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/goph/stdlib/ext"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	config := &configuration{}

	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	// Load configuration from flags
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config.flags(flags)
	flags.Parse(os.Args[1:])

	// Create a new logger
	logger := newLogger(config)
	defer ext.Close(logger)

	// Create a new error handler
	errorHandler := newErrorHandler(config, logger)
	defer ext.Close(errorHandler)

	// Register error handler to recover from panics
	defer emperror.HandleRecover(errorHandler)

	// Application context
	app := &application{
		config:          config,
		logger:          logger,
		errorHandler:    errorHandler,
		healthCollector: healthz.Collector{},
		tracer:          newTracer(config),
	}

	status := healthz.NewStatusChecker(healthz.Healthy)
	app.healthCollector.RegisterChecker(healthz.ReadinessCheck, status)

	level.Info(logger).Log(
		"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
		"version", Version,
		"commit_hash", CommitHash,
		"build_date", BuildDate,
		"environment", config.Environment,
	)

	serverQueue := newServerQueue(app)
	defer serverQueue.Close()

	errChan := serverQueue.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		status.SetStatus(healthz.Unhealthy)
		level.Debug(logger).Log("msg", "error received from error channel")
		emperror.HandleIfErr(errorHandler, err)
	case s := <-signalChan:
		level.Info(logger).Log("msg", fmt.Sprintf("captured %v", s))
		status.SetStatus(healthz.Unhealthy)

		level.Debug(logger).Log(
			"msg", "shutting down with timeout",
			"timeout", config.ShutdownTimeout,
		)

		ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

		err := serverQueue.Shutdown(ctx)
		if err != nil {
			errorHandler.Handle(err)
		}

		// Cancel context if shutdown completed earlier
		cancel()
	}
}
