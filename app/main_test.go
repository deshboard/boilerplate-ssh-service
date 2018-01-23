package app

import (
	"os"
	"testing"

	"github.com/goph/fxt/test"
	"github.com/goph/fxt/test/is"
)

var acceptanceRunner *test.GodogRunner

func TestMain(m *testing.M) {
	result := 0

	var runners []func() int

	if is.Acceptance {
		runners = append(runners, acceptanceRunner.Run)
	}

	if is.Unit || is.Integration || !is.Acceptance {
		runners = append(runners, m.Run)
	}

	for _, run := range runners {
		if r := run(); r > result {
			result = r
		}
	}

	os.Exit(result)
}
