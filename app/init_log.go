package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/logrus_fluent"
	"gopkg.in/airbrake/gobrake.v2"
	logrus_airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v2"
)

func init() {
	// Register shutdown handler in logrus
	logrus.RegisterExitHandler(shutdownManager.Shutdown)

	// Log debug level messages if debug mode is turned on
	if config.Debug {
		logger.Logger.Level = logrus.DebugLevel
	}

	// Initialize Airbrake
	if config.AirbrakeEnabled {
		airbrakeHook := logrus_airbrake.NewHook(config.AirbrakeProjectID, config.AirbrakeAPIKey, config.Environment)
		airbrake := airbrakeHook.Airbrake

		airbrake.SetHost(config.AirbrakeEndpoint)

		airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
			notice.Context["version"] = Version
			notice.Context["commit"] = CommitHash

			return notice
		})

		logger.Logger.Hooks.Add(airbrakeHook)
		shutdownManager.Register(airbrake.Close)
	}

	// Initialize Fluentd
	if config.FluentdEnabled {
		fluentdHook, err := logrus_fluent.New(config.FluentdHost, config.FluentdPort)
		if err != nil {
			logger.Panic(err)
		}

		// Configure fluent tag
		fluentdHook.SetTag(LogTag)

		fluentdHook.AddFilter("error", logrus_fluent.FilterError)

		logger.Logger.Hooks.Add(fluentdHook)
		shutdownManager.Register(fluentdHook.Fluent.Close)
	}
}
