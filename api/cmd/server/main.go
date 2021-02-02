package main

import (
	"fmt"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/db"
	"github.com/cs3305-team-4/api/pkg/routes"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {
	log.Info("grindsapp api starting")
	bindStr := fmt.Sprintf(
		"%s:%s",
		viper.GetString("bind.address"),
		viper.GetString("bind.port"),
	)
	// do an initial db connection to test credentials
	if _, err := db.DB(); err != nil {
		panic(fmt.Errorf("could not connect to database: %s", err))
	}
	log.Infof("binding to %s", bindStr)
	router := routes.GetHandler()
	http.ListenAndServe(bindStr, router)
}
