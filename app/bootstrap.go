package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
)

func bootstrap() serverz.Server {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	checkerCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdownManager.Register(w.Close)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  mux,
			ErrorLog: log.New(w, "http: ", 0),
		},
		Name: "http",
	}
}
