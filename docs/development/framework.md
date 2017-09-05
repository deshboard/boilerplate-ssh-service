# Framework


## Overview

This project does not use any third-party framework (except ones required by the application logic),
but relies heavily on the standard library and separate third-party components.
The integration layer for these components and the main execution logic can be found in the [cmd/](../cmd/) directory.


## Components

### Metrics

In order to effectively use metrics, you need to choose a metrics reporting mechanism.
Common mechanisms are push-based (eg. StatsD) and pull-based (eg. Prometheus).
In order to support both, this project comes with [go-kit's metrics package](https://github.com/go-kit/kit) installed,
which is a metric reporting abstraction.

When using a pull-based solution you most likely have to register an HTTP endpoint in the debug server.
You can do so in [debug.go](../cmd/debug.go)

``` go
package main

import (
	"github.com/goph/stdlib/net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// registerDebugRoutes allows to register custom routes in the debug server.
func registerDebugRoutes(appCtx *application, h http.HandlerAcceptor) {
	h.Handle("/metrics", promhttp.Handler())
}
```
