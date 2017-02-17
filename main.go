package main // import "github.com/deshboard/boilerplate-service"
import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-service/app"
	"github.com/sagikazarmark/healthz"
)

// Global context
var (
	config  = &app.Configuration{}
	log     = logrus.New()
	closers = []io.Closer{}
)

func main() {
	defer shutdown()

	var (
		healthAddr = flag.String("health", "0.0.0.0:90", "Health service address.")
	)
	flag.Parse()

	healthHandler, status := healthService()
	healthServer := &http.Server{
		Addr:    *healthAddr,
		Handler: healthHandler,
	}

	log.Printf("Starting %s service", app.FriendlyServiceName)
	log.Printf("Version %s (%s) built at %s", app.Version, app.CommitHash, app.BuildDate)
	log.Printf("Environment: %s", config.Environment)
	log.Printf("%s Health service listening on %s", app.FriendlyServiceName, healthServer.Addr)

	errChan := make(chan error, 10)

	go func() {
		errChan <- healthServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Panic(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			status.SetStatus(healthz.Unhealthy)
			os.Exit(0)
		}
	}
}

// Panic recovery and close handler
func shutdown() {
	v := recover()
	if v != nil {
		log.Error(v)
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
