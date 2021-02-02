package main

import (
	"fmt"
	"net/http"

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

	log.Infof("binding to %s", bindStr)
	router := routes.GetHandler()
	http.ListenAndServe(bindStr, router)
}
