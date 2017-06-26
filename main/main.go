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
	"github.com/goph/serverz"
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
	logger, closer := newLogger(config)
	defer closer.Close()

	// Create a new error handler
	errorHandler, closer := newErrorHandler(config, logger)
	defer closer.Close()

	// Register error handler to recover from panics
	defer emperror.HandleRecover(errorHandler)

	healthCollector := healthz.Collector{}
	tracer := newTracer(config)
	serverQueue := serverz.NewQueue(&serverz.Manager{Logger: logger})
	signalChan := make(chan os.Signal, 1)

	level.Info(logger).Log(
		"msg", fmt.Sprintf("Starting %s", FriendlyServiceName),
		"version", Version,
		"commitHash", CommitHash,
		"buildDate", BuildDate,
		"environment", config.Environment,
	)

	if config.Debug {
		debugServer := newDebugServer(logger)
		serverQueue.Append(debugServer, config.DebugAddr)
		defer debugServer.Close()
	}

	server, closer := newServer(config, logger, errorHandler, tracer, healthCollector)
	serverQueue.Prepend(server, config.ServiceAddr)
	defer closer.Close()
	defer server.Close()

	healthServer, status := newHealthServer(logger, healthCollector)
	serverQueue.Prepend(healthServer, config.HealthAddr)
	defer healthServer.Close()

	errChan := serverQueue.Start()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			status.SetStatus(healthz.Unhealthy)
			level.Debug(logger).Log("msg", "Error received from error channel")
			emperror.HandleIfErr(errorHandler, err)

			// Break the loop, proceed with regular shutdown
			break MainLoop
		case s := <-signalChan:
			level.Info(logger).Log("msg", fmt.Sprintf("Captured %v", s))
			status.SetStatus(healthz.Unhealthy)

			level.Debug(logger).Log(
				"msg", "Shutting down with timeout",
				"timeout", config.ShutdownTimeout,
			)

			ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

			err := serverQueue.Stop(ctx)
			if err != nil {
				errorHandler.Handle(err)
			}

			// Cancel context if shutdown completed earlier
			cancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}
}
