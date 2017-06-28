package main

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/stdlib/ext"
)

// newLogger creates a new logger instance
func newLogger(config *configuration) (log.Logger, ext.Closer) {
	var logger log.Logger

	// Use JSON when in production
	if "production" == config.Environment {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	}

	// Add default context
	logger = log.With(logger, "service", ServiceName)

	// Default to Info level
	logger = level.NewInjector(logger, level.InfoValue())

	// Only log debug level messages if debug mode is turned on
	if config.Debug == false {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return logger, ext.NoopCloser
}
