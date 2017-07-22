# Framework


## Overview

This project does not use any third-party framework (except ones required by the application logic), but relies heavily on the standard library and separate third-party components. The integration layer for these components and the main execution logic can be found in the [main/](../main/) directory.


## Components

### Metrics

In order to effectively use metrics, you need to choose a metrics reporting mechanism. Common mechanisms are push-based (eg. StatsD) and pull-based (eg. Prometheus). In order to support both, this project comes with [Tally](https://github.com/uber-go/tally) installed, which is a metric reporting abstraction.

All you need to do is chosing a reporter implementation in [metrics.go](../main/metrics.go):

``` go
package main

import (
	"io"

    promreporter "github.com/uber-go/tally/prometheus"
	"github.com/uber-go/tally"
)

// newMetrics returns a new tally.Scope used as a root scope.
func newMetrics(config *configuration) interface {
	tally.Scope
	io.Closer
} {
	options := tally.ScopeOptions{
        CachedReporter: promreporter.NewReporter(promreporter.Options{}),
    }

	scope, closer := tally.NewRootScope(options, MetricReportInterval)

	return struct {
		tally.Scope
		io.Closer
        promreporter.Reporter
	}{
		Scope:  scope,
		Closer: closer,
        Reporter: options.CachedReporter,
	}
}
```

In this case the health server implementation detects that this is a pull-based reporter and automatically exposes it under the `/metrics` endpoint.
