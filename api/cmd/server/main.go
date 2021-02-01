package main

import (
	"fmt"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/db"
	"github.com/cs3305-team-4/api/pkg/routes"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("grindsapp api starting")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Could not read config file %s\n", err))
	}

	bindStr := fmt.Sprintf(
		"%s:%s",
		viper.GetString("bind.address"),
		viper.GetString("bind.port"),
	)

	// do an initial db connection to test credentials
	_, err = db.DB()
	if err != nil {
		panic(fmt.Errorf("could not connect to database: %s", err))
	}

	fmt.Printf("binding to %s\n", bindStr)

	router := routes.GetHandler()
	http.ListenAndServe(bindStr, router)
}
