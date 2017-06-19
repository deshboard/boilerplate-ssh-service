package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	defer shutdownManager.Shutdown()

	flag.Parse()

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
		debugServer := &serverz.NamedServer{
			Server: &http.Server{
				Handler:  http.DefaultServeMux,
				ErrorLog: log.New(logWriter, "debug: ", 0),
			},
			Name: "debug",
		}

		serverQueue.Append(debugServer, config.DebugAddr)
		shutdownManager.RegisterAsFirst(debugServer.Close)
	}

	server := bootstrap()

	serverQueue.Prepend(server, config.ServiceAddr)
	shutdownManager.RegisterAsFirst(server.Close)

	status := healthz.NewStatusChecker(healthz.Healthy)
	checkerCollector.RegisterChecker(healthz.ReadinessCheck, status)

	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", checkerCollector.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", checkerCollector.Handler(healthz.ReadinessCheck))

	if config.MetricsEnabled {
		logger.Debug("Serving metrics under health endpoint")

		healthHandler.Handle("/metrics", promhttp.Handler())
	}

	healthServer := &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: log.New(logWriter, "health: ", 0),
		},
		Name: "health",
	}

	serverQueue.Prepend(healthServer, config.HealthAddr)
	shutdownManager.RegisterAsFirst(healthServer.Close)

	errChan := serverQueue.Start()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			status.SetStatus(healthz.Unhealthy)

			if err != nil {
				logger.Error(err)
			}

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
