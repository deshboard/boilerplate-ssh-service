package main

import (
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/opentracing/opentracing-go"
)

// application collects all dependencies and exposes them in a single service locator.
//
// Although service location is a common anti-pattern, it's only purpose here is bootstrapping
// certain parts of the application. DI would be a more appropriate solution, but even there
// bootstrapping requires a single resolution of all dependencies.
type application struct {
	config          *configuration
	logger          log.Logger
	errorHandler    emperror.Handler
	healthCollector healthz.Collector
	tracer          opentracing.Tracer
}
