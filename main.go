package main // import "github.com/deshboard/boilerplate-service"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sagikazarmark/healthz"
)

var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	var (
		healthAddr = flag.String("health", "0.0.0.0:81", "Health service address.")
	)
	flag.Parse()

	healthHandler, status := healthService()

	healthServer := &http.Server{
		Addr:    *healthAddr,
		Handler: healthHandler,
	}

	log.Println("Starting Boilerplate service...")
	log.Printf("Version %s (%s) built at %s", version, commitHash, buildDate)
	log.Printf("Boilerplate Health service listening on %s", healthServer.Addr)

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
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			status.SetStatus(false)
			os.Exit(0)
		}
	}
}

func healthService() (http.Handler, *healthz.StatusHealthChecker) {
	status := healthz.NewStatusHealthChecker(true)
	readinessProbe := healthz.NewProbe()

	healthService := healthz.NewHealthService(healthz.NewProbe(), readinessProbe)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/healthz", healthService.HealthStatus)
	healthMux.HandleFunc("/readiness", healthService.ReadinessStatus)

	return healthMux, status
}
