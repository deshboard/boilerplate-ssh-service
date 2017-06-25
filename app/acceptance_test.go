// +build acceptance

package app

import (
	"os"
	"time"

	"github.com/DATA-DOG/godog"
)

func init() {
	runs = append(runs, func() int {
		format := "progress"
		for _, arg := range os.Args[1:] {
			// go test transforms -v option
			if arg == "-test.v=true" {
				format = "pretty"
				break
			}
		}

		return godog.RunWithOptions(
			"godog",
			FeatureContext,
			godog.Options{
				Format:    format,
				Paths:     []string{"features"},
				Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
			},
		)
	})
}

func FeatureContext(s *godog.Suite) {
	// Add steps here
}
