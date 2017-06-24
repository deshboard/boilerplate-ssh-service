package main

import (
	"flag"
	"time"
)

// configuration holds any kind of config that is necessary for running
type configuration struct {
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

// flags configures a flagset
//
// Note: the current behaviour relies on the fact that at this point environment variables are already loaded.
func (c *configuration) flags(flags *flag.FlagSet) {
	defaultAddr := ""

	// Listen on loopback interface in development mode
	if c.Environment == "development" {
		defaultAddr = "127.0.0.1"
	}

	// Load flags into configuration
	flags.StringVar(&c.ServiceAddr, "service", defaultAddr+":80", "Service address.")
	flags.StringVar(&c.HealthAddr, "health", defaultAddr+":10000", "Health service address.")
	flags.StringVar(&c.DebugAddr, "debug", defaultAddr+":10001", "Debug service address.")
	flags.DurationVar(&c.ShutdownTimeout, "shutdown", 2*time.Second, "Shutdown timeout.")
}
