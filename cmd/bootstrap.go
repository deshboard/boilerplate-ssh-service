package main

import (
	"github.com/goph/fw"
	"github.com/goph/fw/log"
)

// bootstrap bootstraps the application.
func bootstrap() (*application, error) {
	return newApplication(
		configProvider,
		applicationProvider,
		healthProvider,
	)
}

// applicationProvider provides an fw.Application instance.
func applicationProvider(app *application) error {
	a := fw.NewApplication(
		fw.Logger(log.NewLogger(
			log.FormatString(app.config.LogFormat),
			log.Debug(app.config.Debug),
			log.With(
				"environment", app.config.Environment,
				"service", ServiceName,
				"tag", LogTag,
			),
		)),
	)

	app.Application = a

	return nil
}
