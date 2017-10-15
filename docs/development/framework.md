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
can be found in the [cmd/](../../cmd/) directory.


## Components

### Logging

The project uses [go-kit's log package](https://github.com/go-kit/kit) as logging solution.
It is the implementation of a proposed [standard interface](https://docs.google.com/document/d/1shW9DZJXOeGbG9Mr9Us9MiaPqmlcVatD_D8lrOXRNMU/mobilebasic)
for logging, which didn't make it into the Go standard library after all
(see the links for details why), but became widely adopted pattern in the Go ecosystem.

Compared to other logging solutions ([logrus](https://github.com/sirupsen/logrus),
[zap](https://github.com/uber-go/zap), [log15](https://github.com/inconshreveable/log15) to name a few) the
go-kit interface is concise, still very powerful, even if it comes with a bit of penalty
because of the interface type allocations. As of Go 1.9 this isn't a huge problem though,
thanks to some [compiler optimizations](http://commaok.xyz/post/interface-allocs/).

fxt comes with a [logger constructor](https://github.com/goph/fxt/blob/master/log/logger.go) which by default
logs to the standard output, allows info level messages, unless debug mode is turned on and falls back to info
level if no level is specified.


### Error handling

Go treats errors as values which is a completely different approach compared to other popular languages
(eg. Java or C#) where errors are usually exceptions with a try-catch model for error handling.
This is not necessarily better or worse, just different with different benefits and downsides of course.

One such downside for example is that errors do not carry any stack trace, so it's impossible to tell
what was going on in the original context. Another major problem with this error model is that
there is no global error handler and there is no way to "catch" errors in a higher layer of the program
since errors are returned as values.

This project comes with two components to solve these issues:

- [github.com/pkg/errors](https://github.com/pkg/errors) library which attaches a stack trace to a newly created error
- [github.com/goph/emperror](https://github.com/goph/emperror) which provides a super simple error handler interface (along with some implementations)

These tools makes error handling and debugging much easier.

The fxt package contains an [error handler constructor](https://github.com/goph/fxt/blob/master/errors/handler.go)
as well, which by default logs all errors using the logger described in the previous section.


### Tracing

With microservices/SOA conquering more and more software architectures in the world, normal application tracing
(usually just logs with some sort of correlation ID) is simply not enough anymore.

That's where [distributed tracing](http://microservices.io/patterns/observability/distributed-tracing.html) comes
into the picture.

This project comes with a solution called [OpenTracing](http://opentracing.io/) which is a
*"A vendor-neutral open standard for distributed tracing."* and a [Cloud Native Computing Foundation](https://cncf.io/)
member project. OpenTracing provides a common API for libraries so that they are not tied to any tracing
solution.

The API needs an implementation as well which is usually application and tech stack dependent.
Since our code does not need to know the actual implementation, this detail is only important
on the application level.

The fxt library includes constructors for both the builtin "global" tracer (which defaults to a noop implementation)
and for [Jaeger](https://github.com/jaegertracing) which is a distributed tracing platform created by Uber
(also a hosted by [CNCF](https://cncf.io/)). Jaeger is not hardly coupled to the application,
OpenTracing is, so any solutions implementing it's API can easily be integrated into the application.

The following code integrates the global tracer into the application:

```go
package main

import (
	"github.com/goph/fxt/tracing"
	"go.uber.org/fx"
)

func main() {
    fx.New(
        //...
    
        fx.Provide(tracing.NewTracer),
    )	
}
```

Integrating tracer is a little bit more complicated. It requires some configuration registered in the application:

Add the following to [bootstrap.go](../../cmd/bootstrap.go)

```go
package main

import (
	"github.com/goph/fxt/tracing/jaeger"
	jaegerclient "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// NewJaegerConfig returns a new Jaeger config.
func NewJaegerConfig(config *Config) *jaeger.Config {
	c := jaeger.NewConfig(ServiceName)
	c.JaegerConfig = jaegercfg.Configuration{
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: config.JaegerAddr,
			LogSpans:           config.Debug,
		},
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaegerclient.SamplerTypeConst,
			Param: 1,
		},
	}

	return c
}
```

Then register Jaeger in [main.go](../../cmd/main.go)

```go
package main

import (
	"github.com/goph/fxt/tracing/jaeger"
	"go.uber.org/fx"
)

func main() {
    fx.New(
        //...
    
        NewJaegerConfig,
        jaeger.NewTracer,
    )	
}
```


### Metrics

In order to effectively use metrics, you need to choose a metrics reporting mechanism.
Common mechanisms are push-based (eg. StatsD) and pull-based (eg. Prometheus).
In order to support both, using [go-kit's metrics package](https://github.com/go-kit/kit)
is recommended, which is a metric reporting abstraction.
(Go Kit comes preinstalled because of it's logger abstraction.)

When using a pull-based solution you most likely have to register an HTTP endpoint in
the debug server. The constructor provided by fxt exposes the handler, so you can
register eg. the Prometheus endpoint using an invoke function:


Add the following to [bootstrap.go](../../cmd/bootstrap.go)

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

Then register it as an invoke function in [main.go](../../cmd/main.go)

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
