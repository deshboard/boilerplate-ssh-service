package main

import (
	"github.com/deshboard/boilerplate-ssh-service/app"
)

const FriendlyServiceName = app.FriendlyServiceName

// NewConfig creates the application Config from flags and the environment.
func NewConfig() app.Config {
	return app.Config{}
}

// NewApp creates a new application.
func NewApp(config app.Config) *app.Application {
	info := app.ApplicationInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
	}

	return app.NewApp(config, info)
}
