package main

import (
	"flag"
	"time"
)

// defaultTimeout is used as a default for graceful shutdown timeout.
var defaultTimeout = 15 * time.Second

// Config holds any kind of configuration that comes from the outside world and is necessary for running.
type Config struct {
	// Recommended values are: production, development, staging, release/123, etc
	Environment string `default:"production"`
	Debug       bool   `split_words:"true"`
	LogFormat   string `split_words:"true" default:"json"`

	DebugAddr       string        `ignored:"true"`
	ShutdownTimeout time.Duration `ignored:"true"`
}

// flags configures a flagset.
func (c *Config) flags(flags *flag.FlagSet) {
	flags.StringVar(&c.DebugAddr, "debug.addr", ":10000", "Debug and health check address")
	flags.DurationVar(&c.ShutdownTimeout, "shutdown", defaultTimeout, "Timeout for graceful shutdown")
}
