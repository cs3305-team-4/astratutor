package main

import (
	"fmt"
	"net/http"

	"github.com/cs3305-team-4/signalling/internal/app/signalling"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {
	log.Info("grindsapp websocket signalling starting")
	bindStr := fmt.Sprintf(
		"%s:%s",
		viper.GetString("bind.address"),
		viper.GetString("bind.port"),
	)

	log.Infof("websocket binding to: %s", bindStr)
	r := mux.NewRouter()
	r.HandleFunc("/ws/{id}", signalling.ServeWS)
	srv := &http.Server{
		Handler: r,
		Addr:    bindStr,
	}
	srv.ListenAndServe()
}
