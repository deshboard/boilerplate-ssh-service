package main

import (
	"github.com/deshboard/boilerplate-service/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"github.com/goph/fxt"
	"github.com/goph/nest"
	"go.uber.org/fx"
)

// AppModule is an alias so that the main file does not have to import the app package.
var AppModule = app.Module

// NewConfig creates the application Config from flags and the environment.
func NewConfig() (app.Config, error) {
	configurator := nest.NewConfigurator()
	configurator.SetName(app.FriendlyServiceName)

	var config app.Config

	err := configurator.Load(&config)

	return config, err
}

// NewApplicationInfo provides the application information about itself.
func NewApplicationInfo() fxt.ApplicationInfo {
	return fxt.ApplicationInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
	}
}

// Context is a set of dependencies of the application extracted from the container.
type Context struct {
	fx.In

	Config       app.Config
	Runner       app.Runner
	Logger       log.Logger
	ErrorHandler emperror.Handler
}
