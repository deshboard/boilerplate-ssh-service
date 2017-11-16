package main

import (
	"flag"
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

// Context is a set of dependencies of the application extracted from the container.
type Context struct {
	Config       *Config
	Closer       fxt.Closer
	Logger       log.Logger
	ErrorHandler emperror.Handler

	Status   *healthz.StatusChecker
	DebugErr fxdebug.Err
}

// NewApp creates a new application.
func NewApp(config *Config) *Application {
	context := new(Context)

	return &Application{
		App: fx.New(
			fx.NopLogger,
			fxt.Bootstrap,
			fx.Provide(
				func() *Config {
					return config
				},

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

// Status sets the current health status of the application.
func (a *Application) Status(status healthz.Status) {
	a.context.Status.SetStatus(status)
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
}

// NewConfig creates the application Config from flags and the environment.
func NewConfig(flags *flag.FlagSet) *Config {
	config := new(Config)

	config.Flags(flags)

	return config
}

// NewLoggerConfig creates a logger config for the logger constructor.
func NewLoggerConfig(config *Config) (*fxlog.Config, error) {
	c := fxlog.NewConfig()
	f, err := fxlog.ParseFormat(config.LogFormat)
	if err != nil {
		return nil, err
	}

	c.Format = f
	c.Debug = config.Debug
	c.Context = []interface{}{
		"environment", config.Environment,
		"service", ServiceName,
		"tag", LogTag,
	}

	return c, nil
}

// NewDebugConfig creates a debug config for the debug server constructor.
func NewDebugConfig(config *Config) *fxdebug.Config {
	addr := config.DebugAddr

	// Listen on loopback interface in development mode
	if config.Environment == "development" && addr[0] == ':' {
		addr = "127.0.0.1" + addr
	}

	c := fxdebug.NewConfig(addr)
	c.Debug = config.Debug

	return c
}
