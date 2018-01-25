package app

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/goph/fxt/dev"
	"github.com/goph/fxt/log"
	"github.com/goph/fxt/test/nettest"
	"github.com/goph/nest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	dev.LoadEnvFromFile("../.env.test")
	dev.LoadEnvFromFile("../.env.dist")
}

func newConfig() (Config, error) {
	debugPort, _ := nettest.GetFreePort()

	config := Config{
		Environment: "test",
		LogFormat:   "logfmt",
		DebugAddr:   fmt.Sprintf("127.0.0.1:%d", debugPort),
	}

	configurator := nest.NewConfigurator()
	configurator.SetName(FriendlyServiceName)
	configurator.SetArgs([]string{})

	err := configurator.Load(&config)

	return config, err
}

func TestConfig(t *testing.T) {
	defer func() {
		os.Clearenv()
		dev.LoadEnvFromFile("../.env.test")
		dev.LoadEnvFromFile("../.env.dist")
	}()

	tests := map[string]struct {
		env      map[string]string
		args     []string
		actual   Config
		expected Config
	}{
		"full config": {
			map[string]string{
				"ENVIRONMENT": "test",
				"DEBUG":       "false",
				"LOG_FORMAT":  "logfmt",
			},
			[]string{"service", "--debug-addr", ":10001", "--shutdown-timeout", "10s"},
			Config{},
			Config{
				Environment:     "test",
				Debug:           false,
				LogFormat:       log.LogfmtFormat.String(),
				DebugAddr:       ":10001",
				ShutdownTimeout: 10 * time.Second,
			},
		},
		"defaults": {
			map[string]string{},
			[]string{},
			Config{},
			Config{
				Environment:     "production",
				Debug:           false,
				LogFormat:       log.JsonFormat.String(),
				DebugAddr:       ":10000",
				ShutdownTimeout: 15 * time.Second,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range test.env {
				os.Setenv(key, value)
			}

			configurator := nest.NewConfigurator()
			configurator.SetName(FriendlyServiceName)
			configurator.SetArgs(test.args)

			err := configurator.Load(&test.actual)
			require.NoError(t, err)
			assert.Equal(t, test.expected, test.actual)
		})
	}
}
