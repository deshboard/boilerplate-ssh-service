package app

import (
	"testing"

	fxdebug "github.com/goph/fxt/debug"
	fxlog "github.com/goph/fxt/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoggerConfig(t *testing.T) {
	config := Config{
		Environment: "production",
		LogFormat:   "logfmt",
		Debug:       false,
	}

	expected := fxlog.NewConfig()
	expected.Format = fxlog.LogfmtFormat
	expected.Debug = config.Debug
	expected.Context = []interface{}{
		"environment", config.Environment,
		"service", ServiceName,
		"tag", LogTag,
	}

	actual, err := NewLoggerConfig(config)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNewDebugConfig(t *testing.T) {
	tests := map[string]struct {
		config   Config
		expected *fxdebug.Config
	}{
		"production": {
			Config{
				Environment: "production",
				Debug:       false,
				DebugAddr:   ":10000",
			},
			&fxdebug.Config{
				Network: "tcp",
				Addr:    ":10000",
				Debug:   false,
			},
		},
		"development": {
			Config{
				Environment: "development",
				Debug:       true,
				DebugAddr:   ":10000",
			},
			&fxdebug.Config{
				Network: "tcp",
				Addr:    "127.0.0.1:10000",
				Debug:   true,
			},
		},
		"development_with_interface_specified": {
			Config{
				Environment: "development",
				Debug:       true,
				DebugAddr:   "192.168.0.2:10000",
			},
			&fxdebug.Config{
				Network: "tcp",
				Addr:    "192.168.0.2:10000",
				Debug:   true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewDebugConfig(test.config)
			assert.Equal(t, test.expected, actual)
		})
	}

}
