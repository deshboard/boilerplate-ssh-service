package app

import (
	"os"
	"testing"

	"github.com/goph/fxt/test"
	"github.com/goph/fxt/test/is"
)

var runnerFactoryRegistry test.RunnerFactoryRegistry

func TestMain(m *testing.M) {
	runner, err := runnerFactoryRegistry.CreateRunner()
	if err != nil {
		panic(err)
	}

	if is.Unit || is.Integration || !is.Acceptance {
		runner = test.AppendRunner(runner, m)
	}

	result := runner.Run()

	os.Exit(result)
}
