package main

import (
	"flag"

	"github.com/goph/fxt/debug"
	"github.com/goph/fxt/log"
)

// NewConfig creates the application Config from flags and the environment.
func NewConfig(flags *flag.FlagSet) *Config {
	config := new(Config)

	config.Flags(flags)

	return config
}

// NewLoggerConfig creates a logger config for the logger constructor.
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

// NewDebugConfig creates a debug config for the debug server constructor.
func NewDebugConfig(config *Config) *debug.Config {
	addr := config.DebugAddr

	// Listen on loopback interface in development mode
	if config.Environment == "development" && addr[0] == ':' {
		addr = "127.0.0.1" + addr
	}

	c := debug.NewConfig(addr)
	c.Debug = config.Debug

	return c
}
