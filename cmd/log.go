package main

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
)

// loggerProvider creates a new logger instance and registers it in the application.
func loggerProvider(app *application) error {
	var logger log.Logger
	w := log.NewSyncWriter(os.Stdout)

	switch app.config.LogFormat {
	case "logfmt":
		logger = log.NewLogfmtLogger(w)

	case "json":
		logger = log.NewJSONLogger(w)

	default:
		return emperror.NewWithStackTrace(fmt.Sprintf("unsupported log format: %s", app.config.LogFormat))
	}

	// Add default context
	logger = log.With(logger, "service", ServiceName)

	// Default to Info level
	logger = level.NewInjector(logger, level.InfoValue())

	// Only log debug level messages if debug mode is turned on
	if app.config.Debug == false {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	app.logger = logger

	return nil
}
