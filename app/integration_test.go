// +build integration

package app

import (
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
)

func TestMain(m *testing.M) {
	config := &Configuration{}

	envconfig.MustProcess("", config)

	setUp(config)

	result := m.Run()

	tearDown(config)

	os.Exit(result)
}

// Integration test initialization
func setUp(config *Configuration) {
}

// Cleanup after integration tests
func tearDown(config *Configuration) {
}
