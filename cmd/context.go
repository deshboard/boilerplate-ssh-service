package main

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/fxt"
	fxdebug "github.com/goph/fxt/debug"
	"github.com/goph/healthz"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// Context is a set of dependencies of the application extracted from the container.
type Context struct {
	fx.In

	Config       Config
	Logger       log.Logger
	ErrorHandler emperror.Handler

	Status   *healthz.StatusChecker
	DebugErr fxdebug.Err
}

// Wait waits for the application to finish or exit because of some error.
func (c *Context) Wait(app *fxt.App) {
	select {
	case sig := <-app.Done():
		level.Info(c.Logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-c.DebugErr:
		if err != nil {
			err = errors.Wrap(err, "debug server crashed")
			c.ErrorHandler.Handle(err)
		}
	}

	level.Debug(c.Logger).Log("msg", "setting application status to unhealthy")
	c.Status.SetStatus(healthz.Unhealthy)
}
