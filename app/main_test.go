package app

import (
	"testing"

	"os"

	"github.com/joho/godotenv"
)

var runners []func() int

func TestMain(m *testing.M) {
	_ = godotenv.Load("../.env.test", "../.env.dist")

	result := 0

	runners := append(runners, m.Run)

	for _, run := range runners {
		if r := run(); r > result {
			result = r
		}
	}

	os.Exit(result)
}
