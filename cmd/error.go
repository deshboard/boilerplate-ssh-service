package main

import (
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	emperror_log "github.com/goph/emperror/log"
)

// newErrorHandler creates a new Emperror error handler.
func newErrorHandler(config *configuration, logger log.Logger) emperror.Handler {
	return emperror_log.NewHandler(logger)
}
