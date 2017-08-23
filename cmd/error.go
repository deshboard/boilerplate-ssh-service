package main

import (
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"github.com/goph/stdlib/errors"
)

// newErrorHandler creates a new Emperror error handler.
func newErrorHandler(config *configuration, logger log.Logger) errors.Handler {
	return emperror.NewLogHandler(logger)
}
