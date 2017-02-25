package main // import "github.com/deshboard/boilerplate-service"

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
	"github.com/deshboard/boilerplate-service/app"
	"github.com/sagikazarmark/healthz"
)

func main() {
	defer handleShutdown()

	flag.Parse()

	logger.WithFields(logrus.Fields{
		"version":     app.Version,
		"commitHash":  app.CommitHash,
		"buildDate":   app.BuildDate,
		"environment": config.Environment,
	}).Printf("Starting %s service", app.FriendlyServiceName)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdown = append(shutdown, w.Close)

	errChan := make(chan error, 10)

	healthHandler, status := newHealthServiceHandler()
	healthServiceName := fmt.Sprintf("%s Health service", app.FriendlyServiceName)
	healthServer := &http.Server{
		Addr:     config.HealthAddr,
		Handler:  healthHandler,
		ErrorLog: log.New(w, fmt.Sprintf("%s: ", healthServiceName), 0),
	}

	go startHTTPServer(healthServiceName, healthServer)(errChan)

	if config.Debug {
		debugServiceName := fmt.Sprintf("%s Debug service", app.FriendlyServiceName)
		debugServer := &http.Server{
			Addr:     config.DebugAddr,
			Handler:  http.DefaultServeMux,
			ErrorLog: log.New(w, fmt.Sprintf("%s: ", debugServiceName), 0),
		}

		go startHTTPServer(debugServiceName, debugServer)(errChan)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			if err != nil {
				logger.Error(err)
			} else {
				logger.Warning("Error channel received non-error value")
			}

			// Break the loop, proceed with regular shutdown
			break MainLoop
		case s := <-signalChan:
			logger.Println(fmt.Sprintf("Captured %v", s))
			status.SetStatus(healthz.Unhealthy)

			ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				err := healthServer.Shutdown(ctx)
				if err != nil {
					logger.Error(err)
				}

				wg.Done()
			}()

			wg.Wait()

			// Cancel context if shutdown completed earlier
			cancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}

	close(errChan)
	close(signalChan)
}

// Creates a server starter function which can be called as a goroutine
func startHTTPServer(name string, server *http.Server) func(ch chan<- error) {
	// Force close server connections (if graceful closing fails)
	shutdown = append([]shutdownHandler{server.Close}, shutdown...)

	return func(ch chan<- error) {
		logger.WithField("addr", server.Addr).Infof("%s started", name)
		ch <- server.ListenAndServe()
	}
}

type shutdownHandler func() error

// Wraps a function withot error return type
func shutdownFunc(fn func()) shutdownHandler {
	return func() error {
		fn()
		return nil
	}
}

// Panic recovery and shutdown handler
func handleShutdown() {
	v := recover()
	if v != nil {
		logger.Error(v)
	}

	logger.Info("Shutting down")

	for _, handler := range shutdown {
		err := handler()
		if err != nil {
			logger.Error(err)
		}
	}

	if v != nil {
		panic(v)
	}
}
