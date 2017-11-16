package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newConfig() *Config {
	return &Config{
		LogFormat: "logfmt",
	}
}

func TestNewApp(t *testing.T) {
	config := newConfig()
	app := NewApp(config)

	assert.NoError(t, app.Err())
}
