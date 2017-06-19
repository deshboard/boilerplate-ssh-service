package main

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	_fluent "github.com/goph/log/logrus/hooks/fluent"
	"github.com/sirupsen/logrus"
)

func init() {
	logrusLogger := logrus.New()

	// Register shutdown handler in logrus
	logrus.RegisterExitHandler(shutdownManager.Shutdown)

	// Log debug level messages if debug mode is turned on
	if config.Debug {
		logrusLogger.Level = logrus.DebugLevel
	}

	logger.Logger = logrusLogger.WithField("service", ServiceName)

	logWriter = logger.Logger.(*logrus.Entry).WriterLevel(logrus.ErrorLevel)
	shutdownManager.Register(logWriter.Close)

	// Initialize Fluentd
	if config.FluentEnabled {
		f, _ := fluent.New(fluent.Config{
			FluentHost:   config.FluentHost,
			FluentPort:   config.FluentPort,
			AsyncConnect: true, // In case of AsyncConnect there is no error returned
		})

		fluentHook := &_fluent.Hook{
			Fluent: f,
			Tag:    LogTag,
		}

		logrusLogger.Hooks.Add(fluentHook)
		shutdownManager.Register(fluentHook.Fluent.Close)
	}
}
