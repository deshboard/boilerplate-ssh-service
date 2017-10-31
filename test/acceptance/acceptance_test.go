package acceptance

import (
	"testing"

	"flag"
	"os"
	"strconv"
	"time"

	"github.com/DATA-DOG/godog"
)

func TestMain(m *testing.M) {
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

	result := godog.RunWithOptions(
		"godog",
		FeatureContext,
		godog.Options{
			Format:    format,
			Paths:     []string{"../../features"},
			Randomize: seed,
		},
	)

	os.Exit(result)
}
