# Framework


## Overview

This project does not use any third-party framework (except ones required by the application logic), but relies heavily on the standard library and separate third-party components. The integration layer for these components and the main execution logic can be found in the [main/](../main/) directory.


## Components

### Metrics

In order to effectively use metrics, you need to choose a metrics reporting mechanism. Common mechanisms are push-based (eg. StatsD) and pull-based (eg. Prometheus). In order to support both, this project comes with [Tally](https://github.com/uber-go/tally) installed, which is a metric reporting abstraction.

As a first step, you need to return a reporter implementation in [metrics.go](../main/metrics.go):

``` go
package main

import promreporter "github.com/uber-go/tally/prometheus"

// newMetricsReporter returns one of tally.StatsReporter and tally.CachedStatsReporter.
func newMetricsReporter(config *configuration) interface{} {
	return promreporter.NewReporter(promreporter.Options{})
}
```

In this case the health server implementation detects that this is a pull-based reporter and automatically exposes it under the `/metrics` endpoint.

The next step is to create a root scope and start using Tally:

``` go
	scopeOptions := tally.ScopeOptions{
		Prefix: "my_service",
		Tags:   map[string]string{},
	}

    // We support both types of reporters, so check it here.
	if mReporter, ok := metricsReporter.(tally.CachedStatsReporter); ok {
		scopeOptions.CachedReporter = mReporter
		scopeOptions.Separator = promreporter.DefaultSeparator
	}

	scope, closer := tally.NewRootScope(scopeOptions, 1*time.Second)
```
