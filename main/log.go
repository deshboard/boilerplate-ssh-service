package main

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_logrus "github.com/goph/log/logrus"
	_fluent "github.com/goph/logrus-hooks/fluent"
	"github.com/goph/stdlib/ext"
	"github.com/sirupsen/logrus"
)

// newLogger creates a new logger instance
func newLogger(config *configuration) (log.Logger, ext.Closer) {
	logrusLogger, closers := newLogrus(config)
	var logger log.Logger = &_logrus.Logger{Logger: logrusLogger}

	// Default to Info level
	logger = level.NewInjector(logger, level.InfoValue())

	// Only log debug level messages if debug mode is turned on
	if config.Debug == false {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	logger = log.WithPrefix(logger, "service", ServiceName)

	return logger, closers
}

// newLogrus creates a new logrus logger
func newLogrus(config *configuration) (*logrus.Logger, ext.Closers) {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	closers := ext.Closers{}

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

		logger.Hooks.Add(fluentHook)
		closers = append(closers, fluentHook.Fluent)
	}

	return logger, closers
}
