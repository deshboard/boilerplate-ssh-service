// +build acceptance

package app

import (
	"flag"
	"os"
	"time"

	"github.com/DATA-DOG/godog"
)

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

		return godog.RunWithOptions(
			"godog",
			FeatureContext,
			godog.Options{
				Format:    format,
				Paths:     []string{"features"},
				Randomize: seed,
			},
		)
	})
}
