// +build dev

package main

import (
	"time"

	"github.com/goph/fxt/dev"
)

// Load environment configuration in development environment.
func init() {
	dev.LoadEnvFromFile("../.env")
	dev.LoadEnvFromFile("../.env.dist")

	// Load defaults for info variables
	if Version == "" {
		Version = "<dev>"
	}

	if CommitHash == "" {
		CommitHash = "<dev>"
	}

	if BuildDate == "" {
		BuildDate = time.Now().Format(time.RFC3339)
	}
}
