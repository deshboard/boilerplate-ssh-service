package main

import (
	"github.com/airbrake/gobrake"
	"github.com/goph/emperror"
	"github.com/goph/emperror/airbrake"
)

func init() {
	var handlers []emperror.Handler

	// Initialize Airbrake
	if config.AirbrakeEnabled {
		notifier := gobrake.NewNotifier(config.AirbrakeProjectID, config.AirbrakeAPIKey)

		notifier.SetHost(config.AirbrakeEndpoint)

		notifier.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
			if config.Environment == "development" {
				return nil
			}
			notice.Context["environment"] = config.Environment
			notice.Context["version"] = Version
			notice.Context["commit"] = CommitHash

			return notice
		})

		shutdownManager.Register(notifier.Close)

		handlers = append(
			handlers,
			&airbrake.Handler{
				Notifier:          notifier,
				SendSynchronously: true,
			},
		)
	}

	handlers = append(handlers, emperror.NewLogHandler(logger))

	errorHandler = emperror.NewCompositeHandler(handlers...)
}
