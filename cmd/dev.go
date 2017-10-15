// +build dev

package main

import "github.com/joho/godotenv"

func init() {
	// Load environment configuration in development environment.
	// Only works when running the program from the project root.
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.dist")
}
