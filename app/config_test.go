package app

import (
	"fmt"

	"github.com/goph/fxt/test/nettest"
)

func newConfig() Config {
	debugPort, _ := nettest.GetFreePort()

	return Config{
		Environment: "test",
		LogFormat:   "logfmt",
		DebugAddr:   fmt.Sprintf("127.0.0.1:%d", debugPort),
	}
}
