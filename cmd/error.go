package main

import emperror_log "github.com/goph/emperror/log"

// errorHandlerProvider creates a new Emperror error handler and registers it in the application.
func errorHandlerProvider(app *application) error {
	app.errorHandler = emperror_log.NewHandler(app.logger)

	return nil
}
