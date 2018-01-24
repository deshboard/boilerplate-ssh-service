// +build acceptance

package app

import (
	"github.com/DATA-DOG/godog"
	"github.com/goph/fxt/test"
)

func RegisterSuite(runner *test.GodogRunner) {
	runner.RegisterFeaturePath("../features")
	runner.RegisterFeatureContext(FeatureContext)
}

func FeatureContext(s *godog.Suite) {
	// Add steps here
}
