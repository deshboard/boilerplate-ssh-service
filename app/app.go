package app

import (
	fxdebug "github.com/goph/fxt/debug"
	fxerrors "github.com/goph/fxt/errors"
	fxlog "github.com/goph/fxt/log"
	"go.uber.org/fx"
)

// ApplicationInfo is an optional set of information that can be set by the runtime environment (eg. console application).
type ApplicationInfo struct {
	Version    string
	CommitHash string
	BuildDate  string
}

// Module is the collection of all modules of the application.
var Module = fx.Options(
	fx.Provide(
		// Log and error handling
		NewLoggerConfig,
		fxlog.NewLogger,
		fxerrors.NewHandler,

		// Debug server
		NewDebugConfig,
		fxdebug.NewServer,
		fxdebug.NewHealthCollector,
		fxdebug.NewStatusChecker,
	),
)
