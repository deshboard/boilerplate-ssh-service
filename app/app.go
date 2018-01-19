package app

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/fxt"
	fxdebug "github.com/goph/fxt/debug"
	fxerrors "github.com/goph/fxt/errors"
	fxlog "github.com/goph/fxt/log"
	"github.com/goph/healthz"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// Application wraps fx.App and contains a context.
type Application struct {
	*fx.App

	context *Context

	closer fxt.Closer
}

// ApplicationInfo is an optional set of information that can be set by the runtime environment (eg. console application).
type ApplicationInfo struct {
	Version    string
	CommitHash string
	BuildDate  string
}

// Context is a set of dependencies of the application extracted from the container.
type Context struct {
	Config       Config
	Closer       fxt.Closer
	Logger       log.Logger
	ErrorHandler emperror.Handler

	Status   *healthz.StatusChecker
	DebugErr fxdebug.Err
}

// NewApp creates a new application.
func NewApp(config Config, info ApplicationInfo) *Application {
	context := new(Context)

	return &Application{
		App: fx.New(
			fx.NopLogger,
			fxt.Bootstrap,
			fx.Provide(func() (Config, ApplicationInfo) {
				return config, info
			}),
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
			fx.Extract(context),
		),
		context: context,
	}
}

// Close calls the current closer.
func (a *Application) Close() error {
	return a.context.Closer.Close()
}

// Logger returns the application logger.
func (a *Application) Logger() log.Logger {
	return a.context.Logger
}

// ErrorHandler returns the application error handler.
func (a *Application) ErrorHandler() emperror.Handler {
	return a.context.ErrorHandler
}

// Wait waits for the application to finish or exit because of some error.
func (a *Application) Wait() {
	select {
	case sig := <-a.Done():
		level.Info(a.context.Logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-a.context.DebugErr:
		if err != nil {
			err = errors.Wrap(err, "debug server crashed")
			a.context.ErrorHandler.Handle(err)
		}
	}

	a.context.Status.SetStatus(healthz.Unhealthy)
}
