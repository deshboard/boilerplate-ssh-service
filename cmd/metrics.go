package main

import (
	"io"

	"github.com/uber-go/tally"
)

// newMetrics returns a new tally.Scope used as a root scope.
func newMetrics(config *configuration) interface {
	tally.Scope
	io.Closer
} {
	options := tally.ScopeOptions{}

	scope, closer := tally.NewRootScope(options, MetricsReportInterval)

	return struct {
		tally.Scope
		io.Closer
	}{
		Scope:  scope,
		Closer: closer,
	}
}
