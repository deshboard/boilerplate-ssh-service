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
	"syscall"

	"sync"

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

	logger.Printf("Starting %s service", app.FriendlyServiceName)
	logger.Printf("Version %s (%s) built at %s", app.Version, app.CommitHash, app.BuildDate)
	logger.Printf("Environment: %s", config.Environment)

	w := logger.Writer()
	closers = append(closers, w)
	serverLogger := log.New(w, "", 0)

	healthHandler, status := healthService()
	healthServer := &http.Server{
		Addr:     *healthAddr,
		Handler:  healthHandler,
		ErrorLog: serverLogger,
	}

	errChan := make(chan error, 10)

	startHTTPServer(fmt.Sprintf("%s Health", app.FriendlyServiceName), healthServer, errChan)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			if err != nil {
				logger.Panic(err)
			} else {
				logger.Info("Error channel received non-error")

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
					logger.Panic(err)
				}

				wg.Done()
			}()

			wg.Wait()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}

	logger.Info("Shutting down")
}

// Starts an HTTP server
func startHTTPServer(name string, server *http.Server, ch chan<- error) {
	// Force closing server connections (if graceful closing fails)
	closers = append([]io.Closer{server}, closers...)

	logger.Printf("%s service listening on %s", name, server.Addr)

	go func() {
		ch <- server.ListenAndServe()
	}()
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
