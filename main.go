package main // import "github.com/deshboard/boilerplate-service"

import (
	"context"
	"flag"
	"fmt"
	"io"
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

// Global context
var (
	config  = &app.Configuration{}
	logger  = logrus.New()
	closers = []io.Closer{}
)

func main() {
	defer shutdown()

	var (
		healthAddr = flag.String("health", "0.0.0.0:90", "Health service address.")
	)
	flag.Parse()

	logger.WithFields(logrus.Fields{
		"version":     app.Version,
		"commitHash":  app.CommitHash,
		"buildDate":   app.BuildDate,
		"environment": config.Environment,
		"service":     app.ServiceName,
	}).Printf("Starting %s service", app.FriendlyServiceName)

	w := logger.WriterLevel(logrus.ErrorLevel)
	closers = append(closers, w)
	errChan := make(chan error, 10)

	healthHandler, status := healthService()
	healthServer := &http.Server{
		Addr:     *healthAddr,
		Handler:  healthHandler,
		ErrorLog: log.New(w, fmt.Sprintf("%s Health service: ", app.FriendlyServiceName), 0),
	}

	// Force closing server connections (if graceful closing fails)
	closers = append([]io.Closer{healthServer}, closers...)

	go func() {
		logger.WithField("port", healthServer.Addr).Infof("%s Health service started", app.FriendlyServiceName)
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
				panic(err)
			} else {
				logger.Info("Error channel received non-error value")

				// Break the loop, proceed with shutdown
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

	logger.WithField("service", app.ServiceName).Info("Shutting down")
}

// Panic recovery and close handler
func shutdown() {
	v := recover()
	if v != nil {
		logger.Error(v)
	}

	for _, s := range closers {
		s.Close()
	}

	if v != nil {
		panic(v)
	}
}

// Creates the health service and the status checker
func healthService() (http.Handler, *healthz.StatusChecker) {
	status := healthz.NewStatusChecker(healthz.Healthy)
	healthMux := healthz.NewHealthServiceHandler(healthz.NewCheckers(), status)

	return healthMux, status
}
