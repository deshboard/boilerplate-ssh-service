// +build integration

package app_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	integrationSetUp()

	result := m.Run()

	integrationTearDown()

	os.Exit(result)
}

// Integration test initialization
func integrationSetUp() {
}

// Cleanup after integration tests
func integrationTearDown() {
}
