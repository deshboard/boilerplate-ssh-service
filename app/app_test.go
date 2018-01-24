package app

import (
	"testing"
	"time"

	"github.com/goph/fxt"
	"github.com/goph/fxt/test/fxtest"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func newApplicationInfo() fxt.ApplicationInfo {
	return fxt.ApplicationInfo{
		Version:    "<test>",
		CommitHash: "<test>",
		BuildDate:  time.Now().Format(time.RFC3339),
	}
}

func TestApp(t *testing.T) {
	var runner Runner

	app := fxtest.New(
		t,
		fx.NopLogger,
		fx.Provide(newConfig, newApplicationInfo),
		Module,
		fx.Populate(&runner),
	)

	app.RequireStart()

	go func() {
		// TODO: improve this test
		time.Sleep(10 * time.Millisecond)
		app.RequireStop()
	}()

	err := runner.Run(app)
	require.NoError(t, err)

	app.Close()
}
