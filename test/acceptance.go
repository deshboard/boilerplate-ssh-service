package test

import (
	"flag"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/DATA-DOG/godog"
)

var acceptanceRunner *godogRunner

func init() {
	acceptanceRunner = new(godogRunner)
}

type godogRunner struct {
	featurePaths []string
	featureContexts []func(s *godog.Suite)
}

func (r *godogRunner) registerFeaturePath(featurePath string) {
	if path.IsAbs(featurePath) == false {
		_, filename, _, ok := runtime.Caller(2)
		if !ok {
			panic("cannot determine feature path: no caller information")
		}

		featurePath = path.Clean(path.Join(path.Dir(filename), featurePath))
	}

	r.featurePaths = append(r.featurePaths, featurePath)
}

func (r *godogRunner) registerFeatureContext(ctx func(s *godog.Suite)) {
	r.featureContexts = append(r.featureContexts, ctx)
}

func (r *godogRunner) Run() int {
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

	return godog.RunWithOptions(
		"godog",
		func(s *godog.Suite) {
			for _, featureContext := range r.featureContexts {
				featureContext(s)
			}
		},
		godog.Options{
			Format:    format,
			Paths:     r.featurePaths,
			Randomize: seed,
		},
	)
}

func RegisterFeaturePath(featurePath string) {
	acceptanceRunner.registerFeaturePath(featurePath)
}

func RegisterFeatureContext(ctx func(s *godog.Suite)) {
	acceptanceRunner.registerFeatureContext(ctx)
}
