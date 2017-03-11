// +build integration

package app_test

import (
	"os"
	"testing"

	"github.com/deshboard/boilerplate-service/app"
	"github.com/kelseyhightower/envconfig"
)

var config = &app.Configuration{}

func TestMain(m *testing.M) {
	config = &app.Configuration{}

	envconfig.MustProcess("", config)

	setUp(config)

	result := m.Run()

	tearDown(config)

	os.Exit(result)
}

// Integration test initialization
func setUp(config *app.Configuration) {
}

// Cleanup after integration tests
func tearDown(config *app.Configuration) {
}
