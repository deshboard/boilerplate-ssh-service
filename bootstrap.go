package main

import (
	"net/http"

	"github.com/sagikazarmark/healthz"
)

// Creates the health service handler and the status checker
func newHealthServiceHandler() (http.Handler, *healthz.StatusChecker) {
	status := healthz.NewStatusChecker(healthz.Healthy)
	healthMux := healthz.NewHealthServiceHandler(healthz.NewCheckers(), status)

	return healthMux, status
}
