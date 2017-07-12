package main

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// newLogger creates a new logger instance.
func newLogger(config *configuration) log.Logger {
	var logger log.Logger
	w := log.NewSyncWriter(os.Stdout)

	switch config.LogFormat {
	case "logfmt":
		logger = log.NewLogfmtLogger(w)

	case "json":
		logger = log.NewJSONLogger(w)

	default:
		panic(fmt.Errorf("Unsupported log format: %s", config.LogFormat))
	}

	// Add default context
	logger = log.With(logger, "service", ServiceName)

	// Default to Info level
	logger = level.NewInjector(logger, level.InfoValue())

	// Only log debug level messages if debug mode is turned on
	if config.Debug == false {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return logger
}
