package main

import "github.com/opentracing/opentracing-go"

// tracerProvider creates a new Opentracing Tracer and registers it in the application.
func tracerProvider(app *application) error {
	app.tracer = opentracing.GlobalTracer()

	return nil
}
