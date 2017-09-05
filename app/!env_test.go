package app_test

import (
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load("../.env.test", "../.env.dist")
}
