package main

import (
	"github.com/deshboard/boilerplate-service/app"
	"github.com/goph/nest"
)

type Config = app.Config

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
func NewApplicationInfo() app.ApplicationInfo {
	return app.ApplicationInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
	}
}
