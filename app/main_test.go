package app_test

import (
	"testing"

	"os"

	"github.com/joho/godotenv"
)

var runners []func() int

func init() {
	_ = godotenv.Load("../.env.test", "../.env.dist")
}

func TestMain(m *testing.M) {
	result := 0

	runners := append(runners, m.Run)

	for _, run := range runners {
		if r := run(); r > result {
			result = r
		}
	}

	os.Exit(result)
}
