package main

import (
	"net/http"
)

// Returns parameters from the request
// (decouples the service from the router implementation)
type ParamFetcher func(r *http.Request) map[string]string

type Service struct {
	getParams ParamFetcher
}

// Creates a new service object
func NewService() *Service {
	return &Service{
		getParams: func(r *http.Request) map[string]string {
			return make(map[string]string)
		},
	}
}

// This is where the service implementation goes
