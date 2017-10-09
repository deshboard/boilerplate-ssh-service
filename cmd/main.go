package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/fw"
	"github.com/goph/fw-ext/health"
	"github.com/goph/fw/log"
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
		fw.OptionFunc(func(a *fw.Application) fw.ApplicationOption {
			mux, ok := a.Get("debug_handler")
			if _, ok2 := mux.(*http.ServeMux); !ok || !ok2 {
				fw.Entry("debug_handler", http.NewServeMux())(a)
			}

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
