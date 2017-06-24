package main

import (
	"io"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/goph/log"
	_logrus "github.com/goph/log/logrus"
	_fluent "github.com/goph/logrus-hooks/fluent"
	"github.com/goph/stdlib/ext"
	"github.com/sirupsen/logrus"
)

func newLogger(config *Configuration) (log.Logger, io.Writer, ext.Closer) {
	logrusLogger := logrus.New()
	closers := ext.Closers{}

	// Log debug level messages if debug mode is turned on
	if config.Debug {
		logrusLogger.Level = logrus.DebugLevel
	}

	logger := &_logrus.Logger{
		Logger: logrusLogger.WithField("service", ServiceName),
	}

	logWriter := logger.Logger.(*logrus.Entry).WriterLevel(logrus.ErrorLevel)
	closers = append(closers, logWriter)

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
		closers = append(closers, fluentHook.Fluent)
	}

	return logger, logWriter, closers
}
