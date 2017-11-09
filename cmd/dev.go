// +build dev

package main

import (
	"path"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

// Load environment configuration in development environment.
func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot load environment: no caller information")
	}

	root := path.Clean(path.Join(path.Dir(filename), "../"))

	_ = godotenv.Load(path.Join(root, ".env"))
	_ = godotenv.Load(path.Join(root, ".env.dist"))

	// Load defaults for info variables
	if Version == "" {
		Version = "<unknown>"
	}

	if CommitHash == "" {
		CommitHash = "<unknown>"
	}

	if BuildDate == "" {
		BuildDate = time.Now().Format(time.RFC3339)
	}
}
