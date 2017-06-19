package main

import (
	"log"
	"net/http"
	"time"

	"github.com/goph/healthz"
	"github.com/goph/serverz"
)

func bootstrap() serverz.Server {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	checkerCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  mux,
			ErrorLog: log.New(logWriter, "http: ", 0),
		},
		Name: "http",
	}
}
