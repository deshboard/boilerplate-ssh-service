package main

import (
	"fmt"
	"os"
)

// Tries to find an env var and panics if it's not found
func RequireEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable %s is mandatory", key))
	}

	return value
}

// Tries to find an env var and returns a default if it's not found
func DefaultEnv(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return value
}
