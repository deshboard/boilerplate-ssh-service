package main // import "github.com/deshboard/boilerplate-service"

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	logrus_airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v2"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-service/app"
	"github.com/kelseyhightower/envconfig"
	"github.com/sagikazarmark/healthz"
	"gopkg.in/airbrake/gobrake.v2"
)

const (
	serviceName         = "boilerplate.service"
	friendlyServiceName = "Boilerplate"
)

// Provisioned by ldflags
var (
	version    string
	commitHash string
	buildDate  string
)

var (
	config = &app.Configuration{}
	log    = logrus.New()

	airbrake *gobrake.Notifier
)

func init() {
	err := envconfig.Process("app", config)

	if err != nil {
		log.Fatal(err)
	}

	initLog()
}

func main() {
	if config.AirbrakeEnabled {
		defer airbrake.Close()
		defer airbrake.NotifyOnPanic()
	}

	var (
		healthAddr = flag.String("health", "0.0.0.0:90", "Health service address.")
	)
	flag.Parse()

	healthHandler, status := healthService()
	healthServer := &http.Server{
		Addr:    *healthAddr,
		Handler: healthHandler,
	}

	log.Printf("Starting %s service", friendlyServiceName)
	log.Printf("Version %s (%s) built at %s", version, commitHash, buildDate)
	log.Printf("Environment: %s", config.Environment)
	log.Printf("%s Health service listening on %s", friendlyServiceName, healthServer.Addr)

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
			status.SetStatus(healthz.Unhealthy)
			os.Exit(0)
		}
	}
}

// Creates the health service and the status checker
func healthService() (http.Handler, *healthz.StatusChecker) {
	status := healthz.NewStatusChecker(healthz.Healthy)
	healthMux := healthz.NewHealthServiceHandler(healthz.NewCheckers(), status)

	return healthMux, status
}

// Initializes logger
func initLog() {
	if config.AirbrakeEnabled {
		airbrakeHook := logrus_airbrake.NewHook(config.AirbrakeProjectID, config.AirbrakeAPIKey, config.Environment)

		airbrake = airbrakeHook.Airbrake

		airbrake.SetHost(config.AirbrakeHost)

		airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
			notice.Context["version"] = version

			return notice
		})

		log.Hooks.Add(airbrakeHook)
	}
}
