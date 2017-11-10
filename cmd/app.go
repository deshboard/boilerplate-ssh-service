package main

import (
	"github.com/deshboard/boilerplate-service/pkg/context/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"go.uber.org/dig"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService returns a new service instance.
func NewService(params ServiceParams) *context.Service {
	return context.NewService(
		context.Logger(params.Logger),
		context.ErrorHandler(params.ErrorHandler),
	)
}
