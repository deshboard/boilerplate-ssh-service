package main

import (
	"github.com/airbrake/gobrake"
	"github.com/goph/emperror"
	"github.com/goph/emperror/airbrake"
	"github.com/goph/log"
	"github.com/goph/stdlib/ext"
)

func newErrorHandler(config *configuration, logger log.Logger) (emperror.Handler, ext.Closer) {
	var handlers []emperror.Handler
	closers := ext.Closers{}

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

		closers = append(closers, notifier)

		handlers = append(
			handlers,
			&airbrake.Handler{
				Notifier:          notifier,
				SendSynchronously: true,
			},
		)
	}

	handlers = append(handlers, emperror.NewLogHandler(logger))

	errorHandler := emperror.NewCompositeHandler(handlers...)

	return errorHandler, closers
}
