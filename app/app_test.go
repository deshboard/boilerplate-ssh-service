package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newConfig() Config {
	return Config{
		LogFormat: "logfmt",
	}
}

func TestNewApp(t *testing.T) {
	config := newConfig()
	info := ApplicationInfo{
		Version:    "<test>",
		CommitHash: "<test>",
		BuildDate:  time.Now().Format(time.RFC3339),
	}

	app := NewApp(config, info)

	assert.NoError(t, app.Err())
}
