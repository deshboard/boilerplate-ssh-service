package main

import (
	"io"

	"github.com/uber-go/tally"
)

// newMetricScope returns a new tally.Scope used as a root scope.
func newMetricScope(config *configuration) (tally.Scope, io.Closer) {
	options := tally.ScopeOptions{}

	return tally.NewRootScope(options, MetricReportInterval)
}
