package test

import (
	"path"
	"runtime"

	"github.com/joho/godotenv"
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot load environment: no caller information")
	}

	root := path.Clean(path.Join(path.Dir(filename), "../../"))

	_ = godotenv.Load(path.Join(root, ".env.test"))
	_ = godotenv.Load(path.Join(root, ".env.dist"))
}
