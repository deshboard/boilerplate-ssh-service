package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/goph/fxt/test/fxtest"
	"github.com/goph/fxt/test/nettest"
	"go.uber.org/fx"
)

func newConfig() Config {
	debugPort, _ := nettest.GetFreePort()

	return Config{
		Environment:     "test",
		LogFormat:       "logfmt",
		DebugAddr:       fmt.Sprintf("127.0.0.1:%d", debugPort),
		ShutdownTimeout: 15 * time.Second,
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
		fx.Provide(newConfig, newApplicationInfo),
		Module,
	)

	app.RequireStart().RequireStop()
	app.Close()
}
