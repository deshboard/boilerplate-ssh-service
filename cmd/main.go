package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/fw"
	"github.com/goph/fw-ext/debug"
	"github.com/goph/fw-ext/health"
	"github.com/goph/fw/log"
	"github.com/goph/serverz"
	"github.com/kelseyhightower/envconfig"
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

	app := fw.NewApplication(
		fw.Entry("config", config),
		fw.Logger(log.NewLogger(
			log.FormatString(config.LogFormat),
			log.Debug(config.Debug),
			log.With(
				"environment", config.Environment,
				"service", ServiceName,
				"tag", LogTag,
			),
		)),
		fw.LifecycleHook(fw.SignalHook),
		fw.OptionFunc(health.HealthCollector),
		fw.OptionFunc(health.ApplicationStatus),
		fw.OptionFunc(func(a *fw.Application) fw.ApplicationOption {
			return fw.LifecycleHook(fw.Hook{
				PreStart: func() error {
					level.Info(a.Logger()).Log(
						"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
						"version", Version,
						"commit_hash", CommitHash,
						"build_date", BuildDate,
					)

					return nil
				},
			})
		}),
		debug.DebugServer(serverz.NewAddr("tcp", config.DebugAddr)),
	)
	defer app.Close()

	app.Run()
}
