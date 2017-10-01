package main

import (
	"flag"
	"os"

	"github.com/goph/fw"
	"github.com/goph/healthz"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
)

// application collects all dependencies and exposes them in a single service locator.
//
// Although service location is a common anti-pattern, it's only purpose here is bootstrapping
// certain parts of the application. DI would be a more appropriate solution, but even there
// bootstrapping requires a single resolution of all dependencies.
type application struct {
	*fw.Application
	config          *configuration
	healthCollector healthz.Collector
	tracer          opentracing.Tracer
}

// provider is a mutator for an application registering it's dependencies.
type provider func(app *application) error

// newApplication initializes a new application using the passed providers.
func newApplication(providers ...provider) (*application, error) {
	app := new(application)

	for _, p := range providers {
		err := p(app)
		if err != nil {
			// Returning app, so that already initialized resources can still be closed.
			return app, err
		}
	}

	return app, nil
}

// configProvider registers configuration in the application.
func configProvider(app *application) error {
	config := new(configuration)

	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		return err
	}

	// Load configuration from flags
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config.flags(flags)
	flags.Parse(os.Args[1:])

	app.config = config

	return nil
}

// healthProvider registers the health collector in the application.
func healthProvider(app *application) error {
	app.healthCollector = healthz.Collector{}

	return nil
}
