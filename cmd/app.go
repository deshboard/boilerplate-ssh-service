package main

import (
	"flag"

	"github.com/deshboard/boilerplate-service/app"
)

const FriendlyServiceName = app.FriendlyServiceName

// NewConfig creates the application Config from flags and the environment.
func NewConfig(flags *flag.FlagSet) *app.Config {
	config := new(app.Config)

	config.Flags(flags)

	return config
}

// NewApp creates a new application.
func NewApp(config *app.Config) *app.Application {
	info := &app.ApplicationInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
	}

	return app.NewApp(config, info)
}
