package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/fxt"
	"github.com/goph/fxt/debug"
	"github.com/goph/fxt/errors"
	fxlog "github.com/goph/fxt/log"
	"github.com/goph/healthz"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

func main() {
	config := new(configuration)

	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	// Load configuration from flags
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config.flags(flags)
	flags.Parse(os.Args[1:])

	status := healthz.NewStatusChecker(healthz.Healthy)
	var ext struct {
		Closer       fxt.Closer
		Logger       log.Logger
		ErrorHandler emperror.Handler

		DebugErr debug.Err
	}

	app := fx.New(
		fx.NopLogger,
		fxt.Bootstrap,
		fx.Provide(
			NewLoggerConfig(config),
			fxlog.NewLogger,
			errors.NewHandler,
		),
		fx.Provide(
			NewDebugConfig(config),
			debug.NewServer,
			debug.NewHealthCollector,
		),
		fx.Invoke(func(collector healthz.Collector) {
			collector.RegisterChecker(healthz.ReadinessCheck, status)
		}),
		fx.Extract(&ext),
	)

	// Close resources even when there is an error
	defer ext.Closer.Close()

	// Register error handler to recover from panics
	defer emperror.HandleRecover(ext.ErrorHandler)

	err = app.Err()
	if err != nil {
		panic(err)
	}

	level.Info(ext.Logger).Log(
		"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
		"version", Version,
		"commit_hash", CommitHash,
		"build_date", BuildDate,
	)

	err = app.Start(context.Background())
	if err != nil {
		// Try gracefully stopping already started resources
		app.Stop(context.Background())

		panic(err)
	}

	select {
	case sig := <-app.Done():
		status.SetStatus(healthz.Unhealthy)
		level.Info(ext.Logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-ext.DebugErr:
		status.SetStatus(healthz.Unhealthy)

		if err != nil {
			err = emperror.WithStack(emperror.WithMessage(err, "debug server crashed"))
			ext.ErrorHandler.Handle(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	err = app.Stop(ctx)
	emperror.HandleIfErr(ext.ErrorHandler, err)
}
