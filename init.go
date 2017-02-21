package main

import (
	"github.com/deshboard/boilerplate-service/app"
	"github.com/evalphobia/logrus_fluent"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/airbrake/gobrake.v2"
	logrus_airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v2"
)

func init() {
	err := envconfig.Process("app", config)
	if err != nil {
		logger.Fatal(err)
	}

	// Initialize Airbrake
	if config.AirbrakeEnabled {
		airbrakeHook := logrus_airbrake.NewHook(config.AirbrakeProjectID, config.AirbrakeAPIKey, config.Environment)
		airbrake := airbrakeHook.Airbrake

		airbrake.SetHost(config.AirbrakeEndpoint)

		airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
			notice.Context["version"] = app.Version
			notice.Context["commit"] = app.CommitHash

			return notice
		})

		logger.Logger.Hooks.Add(airbrakeHook)
		shutdown = append(shutdown, airbrake.Close)
	}

	// Initialize Fluentd
	if config.FluentdEnabled {
		fluentdHook, err := logrus_fluent.New(config.FluentdHost, config.FluentdPort)
		if err != nil {
			logger.Panic(err)
		}

		fluentdHook.SetTag(app.ServiceName)
		fluentdHook.AddFilter("error", logrus_fluent.FilterError)

		logger.Logger.Hooks.Add(fluentdHook)
		shutdown = append(shutdown, fluentdHook.Fluent.Close)
	}
}
