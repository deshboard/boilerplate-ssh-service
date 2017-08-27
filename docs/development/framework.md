# Framework


## Overview

This project does not use any third-party framework (except ones required by the application logic), but relies heavily on the standard library and separate third-party components. The integration layer for these components and the main execution logic can be found in the [cmd/](../cmd/) directory.


## Components

### Metrics

In order to effectively use metrics, you need to choose a metrics reporting mechanism. Common mechanisms are push-based (eg. StatsD) and pull-based (eg. Prometheus). In order to support both, this project comes with [Tally](https://github.com/uber-go/tally) installed, which is a metric reporting abstraction.

All you need to do is chosing a reporter implementation in [metrics.go](../cmd/metrics.go):

``` go
package main

import (
	"net/http"

	"github.com/goph/stdlib/ext"
	"github.com/uber-go/tally"
	promreporter "github.com/uber-go/tally/prometheus"
)

// newMetrics returns a new tally.Scope used as a root scope.
func newMetrics(config *configuration) interface {
	tally.Scope
	ext.Closer
} {
	reporter := promreporter.NewReporter(promreporter.Options{})

	options := tally.ScopeOptions{
		CachedReporter: reporter,
	}

	scope, closer := tally.NewRootScope(options, MetricsReportInterval)

	return struct {
		tally.Scope
		ext.Closer
		http.Handler
	}{
		Scope:   scope,
		Closer:  closer,
		Handler: reporter.HTTPHandler(),
	}
}
```

In this case the health server implementation detects that this is a pull-based reporter and automatically exposes it under the `/metrics` endpoint.
