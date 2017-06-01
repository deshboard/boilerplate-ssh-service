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

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
)

func main() {
	defer logger.Info("Shutting down")
	defer shutdownManager.Shutdown()

	flag.Parse()

	logger.WithFields(logrus.Fields{
		"version":     Version,
		"commitHash":  CommitHash,
		"buildDate":   BuildDate,
		"environment": config.Environment,
	}).Infof("Starting %s", FriendlyServiceName)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdownManager.Register(w.Close)

	serverManager := serverz.NewServerManager(logger)
	errChan := make(chan error, 10)
	signalChan := make(chan os.Signal, 1)

	var debugServer serverz.Server
	if config.Debug {
		debugServer = &serverz.NamedServer{
			Server: &http.Server{
				Handler:  http.DefaultServeMux,
				ErrorLog: log.New(w, "debug: ", 0),
			},
			Name: "debug",
		}

		shutdownManager.RegisterAsFirst(debugServer.Close)
		go serverManager.ListenAndStartServer(debugServer, config.DebugAddr)(errChan)
	}

	server := bootstrap()

	status := healthz.NewStatusChecker(healthz.Healthy)
	checkerCollector.RegisterChecker(healthz.ReadinessCheck, status)

	healthService := checkerCollector.NewHealthService()
	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", healthService.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", healthService.Handler(healthz.ReadinessCheck))

	if config.MetricsEnabled {
		logger.Debug("Serving metrics under health endpoint")

		healthHandler.Handle("/metrics", promhttp.Handler())
	}

	healthServer := &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: log.New(w, "health: ", 0),
		},
		Name: "health",
	}

	shutdownManager.RegisterAsFirst(server.Close, healthServer.Close)
	go serverManager.ListenAndStartServer(server, config.ServiceAddr)(errChan)
	go serverManager.ListenAndStartServer(healthServer, config.HealthAddr)(errChan)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			status.SetStatus(healthz.Unhealthy)

			if err != nil {
				logger.Error(err)
			} else {
				logger.Warning("Error channel received non-error value")
			}

			// Break the loop, proceed with regular shutdown
			break MainLoop
		case s := <-signalChan:
			logger.Infof(fmt.Sprintf("Captured %v", s))
			status.SetStatus(healthz.Unhealthy)

			logger.Debugf("Shutting down with '%v' timeout", config.ShutdownTimeout)

			ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
			wg := &sync.WaitGroup{}

			if config.Debug {
				go serverManager.StopServer(debugServer, wg)(ctx)
			}
			go serverManager.StopServer(server, wg)(ctx)
			go serverManager.StopServer(healthServer, wg)(ctx)

			wg.Wait()

			// Cancel context if shutdown completed earlier
			cancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}
}
