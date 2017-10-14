# Framework


## Overview

The project uses Uber's [fx](https://github.com/uber-go/fx) application framework
to manage application lifecycle and to resolve dependencies between application
components.

These components may be:

- part of the standard library
- separate third-party libraries
- components provided by [fxt](https://github.com/goph/fxt),
a set of constructors for common components

The integration layer and the main execution logic
can be found in the [cmd/](../cmd/) directory.


## Components

### Metrics

In order to effectively use metrics, you need to choose a metrics reporting mechanism.
Common mechanisms are push-based (eg. StatsD) and pull-based (eg. Prometheus).
In order to support both, using [go-kit's metrics package](https://github.com/go-kit/kit)
is recommended, which is a metric reporting abstraction.
(Go Kit comes preinstalled because of it's logger abstraction.)

When using a pull-based solution you most likely have to register an HTTP endpoint in
the debug server. The constructor provided by fxt exposes the handler, so you can
register eg. the Prometheus endpoint using an invoke function:


Add the following to [bootstrap.go](../cmd/bootstrap.go)

```go
package main

import (
	"github.com/goph/fxt/debug"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterPrometheusHandler registers the Prometheus metrics handler in the debug server.
func RegisterPrometheusHandler(handler debug.Handler) {
	handler.Handle("/metrics", promhttp.Handler())
}
```

Then register it as an invoke function in [main.go](../cmd/main.go)

```go
package main

import (
	"go.uber.org/fx"
)

func main() {
    fx.New(
        //...
    
        fx.Invoke(RegisterPrometheusHandler),
    )	
}
```
