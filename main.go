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

func main() {
	var (
		healthAddr = flag.String("health", "0.0.0.0:81", "Health service address.")
	)
	flag.Parse()

	status := healthz.NewStatusHealthChecker(true)
	readinessProbe := healthz.NewProbe()

	healthService := healthz.NewHealthService(healthz.NewProbe(), readinessProbe)
	healthMux := http.NewServeMux()
	healthService.RegisterHandlers(healthMux)

	healthServer := &http.Server{
		Addr:    *healthAddr,
		Handler: healthMux,
	}

	log.Println("Starting Boilerplate service...")
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
