// +build dev

package main

import "github.com/joho/godotenv"

func init() {
	// Only works when running the program from the project root.
	_ = godotenv.Load(".env", ".env.dist")
}
