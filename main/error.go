package main

import (
	"github.com/goph/emperror"
	"github.com/goph/log"
)

// newErrorHandler creates a new Emperror error handler.
func newErrorHandler(config *configuration, logger log.Logger) emperror.Handler {
	return emperror.NewLogHandler(logger)
}
