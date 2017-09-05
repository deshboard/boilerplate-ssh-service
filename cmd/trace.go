package main

import "github.com/opentracing/opentracing-go"

// newTracer creates a new Opentracing Tracer.
func newTracer(config *configuration) opentracing.Tracer {
	return opentracing.GlobalTracer()
}
