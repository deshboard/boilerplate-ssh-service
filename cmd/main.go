package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/nest"
)

func main() {
	config := NewConfig()

	configurator := nest.NewConfigurator()
	configurator.SetName(FriendlyServiceName)

	err := configurator.Load(config)
	if err == nest.ErrFlagHelp {
		os.Exit(0)
	} else if err != nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	err = app.Stop(ctx)
	if err != nil {
		app.ErrorHandler().Handle(err)
	}
}
