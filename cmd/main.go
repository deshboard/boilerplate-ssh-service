package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/fw"
	"github.com/goph/fw/log"
	"github.com/goph/healthz"
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
		fw.Entry("health_collector", healthz.Collector{}),
		fw.LifecycleHook(fw.SignalHook),
		fw.OptionFunc(func(a *fw.Application) fw.ApplicationOption {
			healthCollector := a.MustGet("health_collector").(healthz.Collector)

			status := healthz.NewStatusChecker(healthz.Healthy)
			healthCollector.RegisterChecker(healthz.ReadinessCheck, status)

			return fw.LifecycleHook(fw.Hook{
				PreStart: func() error {
					status.SetStatus(healthz.Healthy)

					return nil
				},
				PreShutdown: func() error {
					status.SetStatus(healthz.Unhealthy)

					return nil
				},
			})
		}),
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
		fw.OptionFunc(func(a *fw.Application) fw.ApplicationOption {
			debugServer := newDebugServer(a)

			return fw.LifecycleHook(fw.Hook{
				OnStart: func(ctx context.Context, done chan<- interface{}) error {
					lis, err := net.Listen("tcp", config.DebugAddr)
					if err != nil {
						return err
					}

					go func() {
						done <- debugServer.Serve(lis)
					}()

					return nil
				},
				OnShutdown: func(ctx context.Context) error {
					return debugServer.Shutdown(ctx)
				},
			})
		}),
	)
	defer app.Close()

	app.Run()
}
