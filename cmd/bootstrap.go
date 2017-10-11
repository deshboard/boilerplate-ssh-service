package main

import (
	"github.com/goph/fxt/debug"
	"github.com/goph/fxt/log"
	"github.com/goph/serverz"
)

// NewLoggerConfig creates a logger config constructor.
func NewLoggerConfig(config *configuration) func() (*log.Config, error) {
	return func() (*log.Config, error) {
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
}

// NewDebugConfig creates a debug config constructor.
func NewDebugConfig(config *configuration) func() *debug.Config {
	return func() *debug.Config {
		c := debug.NewConfig(serverz.NewAddr("tcp", config.DebugAddr))
		c.Debug = config.Debug

		return c
	}
}
