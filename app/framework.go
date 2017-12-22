package app

import (
	fxdebug "github.com/goph/fxt/debug"
	fxlog "github.com/goph/fxt/log"
)

// NewLoggerConfig creates a logger config for the logger constructor.
func NewLoggerConfig(config Config) (*fxlog.Config, error) {
	c := fxlog.NewConfig()
	f, err := fxlog.ParseFormat(config.LogFormat)
	if err != nil {
		return nil, err
	}

	c.Format = f
	c.Debug = config.Debug
	c.Context = []interface{}{
		"environment", config.Environment,
		"service", ServiceName,
		"tag", LogTag,
	}

	return c, nil
}

// NewDebugConfig creates a debug config for the debug server constructor.
func NewDebugConfig(config Config) *fxdebug.Config {
	addr := config.DebugAddr

	// Listen on loopback interface in development mode
	if config.Environment == "development" && addr[0] == ':' {
		addr = "127.0.0.1" + addr
	}

	c := fxdebug.NewConfig(addr)
	c.Debug = config.Debug

	return c
}
