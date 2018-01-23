package main

import (
	"fmt"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/fxt"
	fxdebug "github.com/goph/fxt/debug"
	"github.com/goph/healthz"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// Runner is a set of dependencies of the application extracted from the container.
type Runner struct {
	fx.In

	Status   *healthz.StatusChecker
	DebugErr fxdebug.Err
}

// Run waits for the application to finish or exit because of some error.
func (r *Runner) Run(app *fxt.App, ctx *Context) {
	select {
	case sig := <-app.Done():
		level.Info(ctx.Logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-r.DebugErr:
		if err != nil {
			err = errors.Wrap(err, "debug server crashed")
			ctx.ErrorHandler.Handle(err)
		}
	}

	level.Debug(ctx.Logger).Log("msg", "setting application status to unhealthy")
	r.Status.SetStatus(healthz.Unhealthy)
}
