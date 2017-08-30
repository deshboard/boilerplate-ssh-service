package app

import (
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
)

// ServiceOption sets options in the Service.
type ServiceOption func(s *Service)

// Logger returns a ServiceOption that sets the logger for the service.
func Logger(l log.Logger) ServiceOption {
	return func(s *Service) {
		s.logger = l
	}
}

// ErrorHandler returns a ServiceOption that sets the error handler for the service.
func ErrorHandler(l emperror.Handler) ServiceOption {
	return func(s *Service) {
		s.errorHandler = l
	}
}

// Service contains the main controller logic.
type Service struct {
	logger       log.Logger
	errorHandler emperror.Handler
}

// NewService creates a new service object.
func NewService(opts ...ServiceOption) *Service {
	s := new(Service)

	for _, opt := range opts {
		opt(s)
	}

	// Default logger
	if s.logger == nil {
		s.logger = log.NewNopLogger()
	}

	// Default error handler
	if s.errorHandler == nil {
		s.errorHandler = emperror.NewNopHandler()
	}

	return s
}
