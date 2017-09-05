package main

import (
	"flag"
	"time"
)

// configuration holds any kind of config that is necessary for running.
type configuration struct {
	// Recommended values are: production, development, staging, release/123, etc
	Environment string `default:"production"`
	Debug       bool   `split_words:"true"`
	LogFormat   string `split_words:"true" default:"json"`

	SSHAddr         string        `ignored:"true"`
	DebugAddr       string        `ignored:"true"`
	ShutdownTimeout time.Duration `ignored:"true"`

	HostPrivateKey         string `split_words:"true"`
	HostPrivateKeyFile     string `split_words:"true"`
	RootLoginAllowed       bool   `split_words:"true"`
	RootAuthorizedKeys     string `split_words:"true"`
	RootAuthorizedKeysFile string `split_words:"true"`
}

// flags configures a flagset.
//
// Note: the current behaviour relies on the fact that at this point environment variables are already loaded.
func (c *configuration) flags(flags *flag.FlagSet) {
	defaultAddr := ""

	// Listen on loopback interface in development mode
	if c.Environment == "development" {
		defaultAddr = "127.0.0.1"
	}

	// Load flags into configuration
	flags.StringVar(&c.SSHAddr, "ssh.addr", defaultAddr+":2222", "SSH service address.")
	flags.StringVar(&c.DebugAddr, "debug.addr", defaultAddr+":10000", "Debug and health check address.")
	flags.DurationVar(&c.ShutdownTimeout, "shutdown", 2*time.Second, "Shutdown timeout.")
}
