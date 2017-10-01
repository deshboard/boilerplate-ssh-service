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
)

func main() {
	app, err := newApplication(
		configProvider,
		applicationProvider,
		healthProvider,
	)
	// Close resources even when there is an error
	defer app.Close()

	if err != nil {
		// Handle the error and exit if we have an error handler
		if app.ErrorHandler() != nil {
			app.ErrorHandler().Handle(err)

			os.Exit(1)
		}

		// Otherwise panic
		panic(err)
	}

	// Register error handler to recover from panics
	defer emperror.HandleRecover(app.ErrorHandler())

	status := healthz.NewStatusChecker(healthz.Healthy)
	app.healthCollector.RegisterChecker(healthz.ReadinessCheck, status)

	level.Info(app.Logger()).Log(
		"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
		"version", Version,
		"commit_hash", CommitHash,
		"build_date", BuildDate,
	)

	serverQueue := newServerQueue(app)
	defer serverQueue.Close()

	errChan := serverQueue.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		status.SetStatus(healthz.Unhealthy)
		level.Debug(app.Logger()).Log("msg", "error received from error channel")
		emperror.HandleIfErr(app.ErrorHandler(), err)
	case s := <-signalChan:
		level.Info(app.Logger()).Log("msg", fmt.Sprintf("captured %v", s))
		status.SetStatus(healthz.Unhealthy)

		level.Debug(app.Logger()).Log(
			"msg", "shutting down with timeout",
			"timeout", app.config.ShutdownTimeout,
		)

		ctx, cancel := context.WithTimeout(context.Background(), app.config.ShutdownTimeout)

		err := serverQueue.Shutdown(ctx)
		if err != nil {
			app.ErrorHandler().Handle(err)
		}

		// Cancel context if shutdown completed earlier
		cancel()
	}
}
