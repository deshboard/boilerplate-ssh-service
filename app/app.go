package app

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	fxdebug "github.com/goph/fxt/debug"
	fxerrors "github.com/goph/fxt/errors"
	fxlog "github.com/goph/fxt/log"
	"github.com/goph/healthz"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// Module is the collection of all modules of the application.
var Module = fx.Options(
	fx.Provide(
		// Log and error handling
		NewLoggerConfig,
		fxlog.NewLogger,
		fxerrors.NewHandler,

		// Debug server
		NewDebugConfig,
		fxdebug.NewServer,
		fxdebug.NewHealthCollector,
		fxdebug.NewStatusChecker,
	),
)

// Runner executes the application and waits for it to end.
type Runner struct {
	fx.In

	Logger log.Logger
	Status *healthz.StatusChecker

	DebugErr fxdebug.Err
}

// Run waits for the application to finish or exit because of some error.
func (r *Runner) Run(app interface {
	Done() <-chan os.Signal
}) error {
	defer func() {
		level.Debug(r.Logger).Log("msg", "setting application status to unhealthy")
		r.Status.SetStatus(healthz.Unhealthy)
	}()

	select {
	case sig := <-app.Done():
		fmt.Println() // empty line before log entry
		level.Info(r.Logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-r.DebugErr:
		if err != nil && err != fxdebug.ErrServerClosed {
			return errors.Wrap(err, "debug server crashed")
		}
	}

	return nil
}
