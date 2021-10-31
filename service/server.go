package service

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// runServer runs main service http server
func (b *BDayService) runServer() {

	log.WithField("context", "Run").Infof("starting http server on port %s", b.HttpPort)

	serverMux := mux.NewRouter()
	serverMux.HandleFunc("/hello/{username:[a-zA-Z]+}", b.username)

	server := &http.Server{
		Addr:    ":" + b.HttpPort,
		Handler: serverMux,
	}
	b.isReady = true
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("context", "runServer").WithError(err).Fatalf("server listen and serve")
		}
	}()
	b.waitForExitSignal(server, "runServer")
}
