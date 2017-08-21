// +build acceptance

package app_test

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DATA-DOG/godog"
)

var FeatureContext func(s *godog.Suite)

func init() {
	runners = append(runners, func() int {
		format := "progress"
		seed := int64(0)

		// go test transforms -v option
		if verbose := flag.Lookup("test.v"); verbose != nil {
			format = "pretty"
		}

		// Randomize scenario execution order
		if randomize, _ := strconv.ParseBool(os.Getenv("TEST_RANDOMIZE")); randomize {
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
