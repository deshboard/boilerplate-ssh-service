package main

import (
	"flag"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-service/app"
	"github.com/evalphobia/logrus_fluent"
	"github.com/kelseyhightower/envconfig"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/trace"
	"gopkg.in/airbrake/gobrake.v2"
	logrus_airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v2"
)

// Global context variables
var (
	config   = &app.Configuration{}
	logger   = logrus.New().WithField("service", app.ServiceName) // Use logrus.FieldLogger type
	tracer   = opentracing.GlobalTracer()
	shutdown = []shutdownHandler{}
)

func init() {
	// Load configuration from environment
	err := envconfig.Process("app", config)
	if err != nil {
		logger.Fatal(err)
	}

	// Load flags into configuration
	flag.StringVar(&config.ServiceAddr, "service", "127.0.0.1:80", "Service address.")
	flag.StringVar(&config.HealthAddr, "health", "127.0.0.1:10000", "Health service address.")
	flag.StringVar(&config.DebugAddr, "debug", "127.0.0.1:10001", "Debug service address.")
	flag.DurationVar(&config.ShutdownTimeout, "shutdown", 2*time.Second, "Shutdown timeout.")

	// This is probably OK as the service runs in Docker
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}

	// Initialize Airbrake
	if config.AirbrakeEnabled {
		airbrakeHook := logrus_airbrake.NewHook(config.AirbrakeProjectID, config.AirbrakeAPIKey, config.Environment)
		airbrake := airbrakeHook.Airbrake

		airbrake.SetHost(config.AirbrakeEndpoint)

		airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
			notice.Context["version"] = app.Version
			notice.Context["commit"] = app.CommitHash

			return notice
		})

		logger.Logger.Hooks.Add(airbrakeHook)
		shutdown = append(shutdown, airbrake.Close)
	}

	// Initialize Fluentd
	if config.FluentdEnabled {
		fluentdHook, err := logrus_fluent.New(config.FluentdHost, config.FluentdPort)
		if err != nil {
			logger.Panic(err)
		}

		fluentdHook.SetTag(app.ServiceName)
		fluentdHook.AddFilter("error", logrus_fluent.FilterError)

		logger.Logger.Hooks.Add(fluentdHook)
		shutdown = append(shutdown, fluentdHook.Fluent.Close)
	}
}
