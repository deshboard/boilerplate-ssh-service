package main

import (
	"github.com/deshboard/boilerplate-service/app"
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

		airbrake.SetHost(config.AirbrakeHost)

		airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
			notice.Context["version"] = app.Version

			return notice
		})

		logger.Hooks.Add(airbrakeHook)
		closers = append(closers, airbrake)
	}
}
