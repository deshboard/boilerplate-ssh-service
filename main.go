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
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sagikazarmark/healthz"
)

// Global context variables
var (
	config   = &app.Configuration{}
	logger   = logrus.New().WithField("service", app.ServiceName) // Use logrus.FieldLogger type
	tracer   = opentracing.GlobalTracer()
	shutdown = []shutdownHandler{}
)

func main() {
	defer handleShutdown()

	var (
		healthAddr = flag.String("health", "0.0.0.0:90", "Health service address.")
	)
	flag.Parse()

	logger.WithFields(logrus.Fields{
		"version":     app.Version,
		"commitHash":  app.CommitHash,
		"buildDate":   app.BuildDate,
		"environment": config.Environment,
	}).Printf("Starting %s service", app.FriendlyServiceName)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdown = append(shutdown, w.Close)

	healthHandler, status := newHealthServiceHandler()
	healthServer := &http.Server{
		Addr:     *healthAddr,
		Handler:  healthHandler,
		ErrorLog: log.New(w, fmt.Sprintf("%s Health service: ", app.FriendlyServiceName), 0),
	}

	// Force closing server connections (if graceful closing fails)
	shutdown = append([]shutdownHandler{healthServer.Close}, shutdown...)

	errChan := make(chan error, 10)

	go func() {
		logger.WithField("addr", healthServer.Addr).Infof("%s Health service started", app.FriendlyServiceName)
		errChan <- healthServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			// In theory this can only be non-nil
			if err != nil {
				// This will be handled (logged) by shutdown
				panic(err)
			} else {
				logger.Info("Error channel received non-error value")

				// Break the loop, proceed with regular shutdown
				break MainLoop
			}
		case s := <-signalChan:
			logger.Println(fmt.Sprintf("Captured %v", s))
			status.SetStatus(healthz.Unhealthy)

			shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
			defer shutdownCancel()

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				err := healthServer.Shutdown(shutdownContext)
				if err != nil {
					logger.Error(err)
				}

				wg.Done()
			}()

			wg.Wait()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}

	close(errChan)
	close(signalChan)
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
