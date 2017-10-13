package main

import (
	"flag"
	"os"

	"github.com/goph/fxt/debug"
	"github.com/goph/fxt/log"
	"github.com/kelseyhightower/envconfig"
)

// NewConfig creates the application Config from flags and the environment.
func NewConfig() (*Config, error) {
	config := new(Config)

	// Load Config from environment
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	// Load Config from flags
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config.flags(flags)
	flags.Parse(os.Args[1:])

	return config, nil
}

// NewLoggerConfig creates a logger config constructor.
func NewLoggerConfig(config *Config) (*log.Config, error) {
	c := log.NewConfig()
	f, err := log.ParseFormat(config.LogFormat)
	if err != nil {
		return nil, err
	}

	c.Format = f
	c.Debug = config.Debug
	c.Context = []interface{}{
		"environment", config.Environment,
		"service", ServiceName,
		"tag", LogTag,
	}

	return c, nil
}

// NewDebugConfig creates a debug config.
func NewDebugConfig(config *Config) *debug.Config {
	c := debug.NewConfig(config.DebugAddr)
	c.Debug = config.Debug

	return c
}
