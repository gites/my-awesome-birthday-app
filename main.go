package main

import (
	"github.com/gites/my-awesome-birthday-app/service"
	"github.com/golang-migrate/migrate/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	svc, err := service.NewBDayService()
	if err != nil {
		log.WithField("context", "main").WithError(err).Fatalf("couldn't initialize")
	}
	if err := svc.ConnectDB(); err != nil {
		log.WithField("context", "main").WithError(err).Fatalf("couldn't connect to db")
	}
	if err := service.RunMigrations(svc); err != nil && err != migrate.ErrNoChange {
		log.WithField("context", "main").WithError(err).Fatalf("couldn't run db migrations")
	}
	svc.Run()
}
