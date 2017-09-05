package app_test

import (
	"testing"

	"os"
)

var runners []func() int

func TestMain(m *testing.M) {
	result := 0

	runners := append(runners, m.Run)

	for _, run := range runners {
		if r := run(); r > result {
			result = r
		}
	}

	os.Exit(result)
}
