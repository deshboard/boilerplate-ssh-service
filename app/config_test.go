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
	os.Clearenv()
	defer func() {
		os.Clearenv()
		dev.LoadEnvFromFile("../.env.test")
		dev.LoadEnvFromFile("../.env.dist")
	}()

	env := map[string]string{
		"ENVIRONMENT": "test",
		"DEBUG":       "false",
		"LOG_FORMAT":  "json",
	}

	for key, value := range env {
		os.Setenv(key, value)
	}

	expected := Config{
		Environment:     "test",
		Debug:           false,
		LogFormat:       log.JsonFormat.String(),
		DebugAddr:       ":10000",
		ShutdownTimeout: 15 * time.Second,
	}
	actual := Config{}

	configurator := nest.NewConfigurator()
	configurator.SetName(FriendlyServiceName)
	configurator.SetArgs([]string{})

	err := configurator.Load(&actual)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
