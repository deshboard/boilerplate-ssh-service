// +build acceptance

package app_test

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/DATA-DOG/godog"
)

var FeatureContext func(s *godog.Suite)

func init() {
	runners = append(runners, func() int {
		format := "progress"
		seed := int64(0)

		var verbose, randomize bool
		flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		// go test transforms -v option
		flags.BoolVar(&verbose, "test.v", false, "Test verbosity")
		flags.BoolVar(&randomize, "randomize", false, "Randomize acceptance test order")
		flags.Parse(os.Args[1:])

		if verbose {
			format = "pretty"
		}

		// Randomize scenario execution order
		if randomize {
			seed = time.Now().UTC().UnixNano()
		}

		featureContext := FeatureContext
		if featureContext == nil {
			featureContext = func(s *godog.Suite) {
				fmt.Println("No feature context")
			}
		}

		return godog.RunWithOptions(
			"godog",
			featureContext,
			godog.Options{
				Format:    format,
				Paths:     []string{"../features"},
				Randomize: seed,
			},
		)
	})
}
