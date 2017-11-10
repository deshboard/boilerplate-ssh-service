package test

import (
	"os"
	"testing"
)

var unit bool
var acceptance bool
var integration bool

func Main(m *testing.M) {
	result := 0

	var runners []func() int

	if acceptance {
		runners = append(runners, acceptanceRunner.Run)
	}

	if unit || integration || !acceptance {
		runners = append(runners, m.Run)
	}

	for _, run := range runners {
		if r := run(); r > result {
			result = r
		}
	}

	os.Exit(result)
}
