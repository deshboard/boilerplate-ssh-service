package main

import (
	"flag"
	"time"
)

// Configuration holds any kind of config that is necessary for running.
type Configuration struct {
	// Recommended values are: production, development, staging, release/123, etc
	Environment string `default:"production"`
	Debug       bool   `split_words:"true"`

	ServiceAddr     string        `ignored:"true"`
	HealthAddr      string        `ignored:"true"`
	DebugAddr       string        `ignored:"true"`
	ShutdownTimeout time.Duration `ignored:"true"`

	AirbrakeEnabled   bool   `split_words:"true"`
	AirbrakeEndpoint  string `split_words:"true"`
	AirbrakeProjectID int64  `envconfig:"airbrake_project_id"`
	AirbrakeAPIKey    string `envconfig:"airbrake_api_key"`

	FluentEnabled bool   `split_words:"true"`
	FluentHost    string `split_words:"true"`
	FluentPort    int    `split_words:"true" default:"24224"`
}

func configureFlags(config *Configuration, flags *flag.FlagSet) {
	defaultAddr := ""

	// Listen on loopback interface in development mode.
	if config.Environment == "development" {
		defaultAddr = "127.0.0.1"
	}

	// Load flags into configuration.
	flags.StringVar(&config.ServiceAddr, "service", defaultAddr+":80", "Service address.")
	flags.StringVar(&config.HealthAddr, "health", defaultAddr+":10000", "Health service address.")
	flags.StringVar(&config.DebugAddr, "debug", defaultAddr+":10001", "Debug service address.")
	flags.DurationVar(&config.ShutdownTimeout, "shutdown", 2*time.Second, "Shutdown timeout.")
}
