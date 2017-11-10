// +build acceptance

package app_test

import (
	"github.com/DATA-DOG/godog"
	"github.com/deshboard/boilerplate-ssh-service/test"
)

func init() {
	test.RegisterFeaturePath("../features")
	test.RegisterFeatureContext(FeatureContext)
}

func FeatureContext(s *godog.Suite) {
	// Add steps here
}
