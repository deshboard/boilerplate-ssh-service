package app

import (
	"testing"

	"os"
)

var runs []func() int

func TestMain(m *testing.M) {
	result := 0

	runs := append(runs, m.Run)

	for _, run := range runs {
		if r := run(); r > result {
			result = r
		}
	}

	os.Exit(result)
}
