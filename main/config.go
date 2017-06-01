package main

import "time"

// Configuration holds any kind of config that is necessary for running
type Configuration struct {
	// Recommended values are: production, development, staging, release/123, etc
	Environment string `default:"production"`
	Debug       bool   `split_words:"true"`

	ServiceAddr     string        `ignored:"true"`
	HealthAddr      string        `ignored:"true"`
	DebugAddr       string        `ignored:"true"`
	ShutdownTimeout time.Duration `ignored:"true"`

	// Enable Prometheus metrics to be exposed in the health service under /metrics endpoint.
	MetricsEnabled bool `split_words:"true"`

	AirbrakeEnabled   bool   `split_words:"true"`
	AirbrakeEndpoint  string `split_words:"true"`
	AirbrakeProjectID int64  `envconfig:"airbrake_project_id"`
	AirbrakeAPIKey    string `envconfig:"airbrake_api_key"`

	FluentdEnabled bool   `split_words:"true"`
	FluentdHost    string `split_words:"true"`
	FluentdPort    int    `split_words:"true" default:"24224"`
}
