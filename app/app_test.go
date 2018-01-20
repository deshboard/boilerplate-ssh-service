package app

import (
	"testing"
	"time"

	"github.com/goph/fxt"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func newConfig() Config {
	return Config{
		LogFormat: "logfmt",
	}
}

func newApplicationInfo() ApplicationInfo {
	return ApplicationInfo{
		Version:    "<test>",
		CommitHash: "<test>",
		BuildDate:  time.Now().Format(time.RFC3339),
	}
}

func TestApp(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxt.Bootstrap,
		fx.Provide(newConfig, newApplicationInfo),
		Module,
	)

	app.RequireStart().RequireStop()
}
