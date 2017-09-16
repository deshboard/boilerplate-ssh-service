package main

import (
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
)

// newTracer creates a new Opentracing Tracer.
func newTracer(config *configuration, logger log.Logger) opentracing.Tracer {
	return opentracing.GlobalTracer()
}
