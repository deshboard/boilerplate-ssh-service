package main

import (
	"flag"
	"os"
	"testing"
)

var integration = flag.Bool("integration", false, "run integration tests")

func TestMain(m *testing.M) {
	flag.Parse()

	if *integration {
		integrationSetUp()
	}

	result := m.Run()

	if *integration {
		integrationTearDown()
	}

	os.Exit(result)
}
