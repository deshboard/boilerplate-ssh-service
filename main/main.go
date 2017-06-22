package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config := newConfigWithFlags(flags)

	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	// Load configuration from flags
	flags.Parse(os.Args[1:])

	logger, logWriter, closer := newLogger(config)
	defer closer.Close()

	errorHandler, closer := newErrorHandler(config, logger)
	defer closer.Close()
	defer emperror.HandleRecover(errorHandler)

	logger.Info(
		fmt.Sprintf("Starting %s", FriendlyServiceName),
		map[string]interface{}{
			"version":     Version,
			"commitHash":  CommitHash,
			"buildDate":   BuildDate,
			"environment": config.Environment,
		},
	)

	serverManager := &serverz.Manager{
		Logger:       logger,
		ErrorHandler: errorHandler,
	}
	serverQueue := serverz.NewQueue(serverManager)
	signalChan := make(chan os.Signal, 1)

	if config.Debug {
		debugServer := newDebugServer(logWriter)
		serverQueue.Append(debugServer, config.DebugAddr)
		defer debugServer.Close()
	}

	checkerCollector := healthz.Collector{}

	server := bootstrap(config, logWriter, checkerCollector)
	serverQueue.Prepend(server, config.ServiceAddr)
	defer server.Close()

	status := healthz.NewStatusChecker(healthz.Healthy)
	checkerCollector.RegisterChecker(healthz.ReadinessCheck, status)

	healthServer := newHealthServer(logWriter, checkerCollector)
	serverQueue.Prepend(healthServer, config.HealthAddr)
	defer healthServer.Close()

	errChan := serverQueue.Start()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			status.SetStatus(healthz.Unhealthy)
			logger.Debug("Error received from error channel")
			emperror.HandleIfErr(errorHandler, err)

			// Break the loop, proceed with regular shutdown
			break MainLoop
		case s := <-signalChan:
			logger.Info(fmt.Sprintf("Captured %v", s))
			status.SetStatus(healthz.Unhealthy)

			logger.Debug("Shutting down with timeout", map[string]interface{}{"timeout": config.ShutdownTimeout})

			ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
			wg := &sync.WaitGroup{}

			serverQueue.Stop(ctx)

			wg.Wait()

			// Cancel context if shutdown completed earlier
			cancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}
}
