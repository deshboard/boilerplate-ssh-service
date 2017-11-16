package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/kelseyhightower/envconfig"
)

// defaultShutdownTimeout is used as a default for graceful shutdown timeout.
var defaultShutdownTimeout = 15 * time.Second

func main() {
	prefix := flag.String("prefix", "", "Environment variable prefix (useful when multiple apps use the same environment)")
	shutdownTimeout := flag.Duration("shutdown.timeout", defaultShutdownTimeout, "Timeout for graceful shutdown")

	config := NewConfig(flag.CommandLine)

	flag.Parse()

	// Load config from environment (from the appropriate prefix)
	err := envconfig.Process(*prefix, config)
	if err != nil {
		panic(err)
	}

	app := NewApp(config)

	err = app.Err()
	if err != nil {
		panic(err)
	}

	// Close resources when the application stops running
	defer app.Close()

	// Register error handler to recover from panics
	defer emperror.HandleRecover(app.ErrorHandler())

	level.Info(app.Logger()).Log(
		"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
		"version", Version,
		"commit_hash", CommitHash,
		"build_date", BuildDate,
	)

	err = app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	app.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()

	err = app.Stop(ctx)
	if err != nil {
		app.ErrorHandler().Handle(err)
	}
}
