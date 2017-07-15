package main

import (
	"io"

	"github.com/uber-go/tally"
)

// newMetrics returns a new tally.Scope used as a root scope.
func newMetrics(config *configuration) (tally.Scope, io.Closer) {
	options := tally.ScopeOptions{}

	return tally.NewRootScope(options, MetricReportInterval)
}
