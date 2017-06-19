package main

import (
	"flag"
	"io"
	"time"

	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/goph/log/logrus"
	"github.com/goph/shutdown"
	"github.com/kelseyhightower/envconfig"
	opentracing "github.com/opentracing/opentracing-go"
)

// Global context variables
var (
	config           = &Configuration{}
	logger           = &logrus.Logger{}
	logWriter        io.WriteCloser
	errorHandler     emperror.Handler
	tracer           = opentracing.GlobalTracer()
	shutdownManager  = shutdown.NewManager()
	checkerCollector = healthz.Collector{}
)

func init() {
	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	defaultAddr := ""

	// Listen on loopback interface in development mode
	if config.Environment == "development" {
		defaultAddr = "127.0.0.1"
	}

	// Load flags into configuration
	flag.StringVar(&config.ServiceAddr, "service", defaultAddr+":80", "Service address.")
	flag.StringVar(&config.HealthAddr, "health", defaultAddr+":10000", "Health service address.")
	flag.StringVar(&config.DebugAddr, "debug", defaultAddr+":10001", "Debug service address.")
	flag.DurationVar(&config.ShutdownTimeout, "shutdown", 2*time.Second, "Shutdown timeout.")
}
