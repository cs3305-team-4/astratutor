package services

import (
	"fmt"
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Could not read config file %s\n", err))
	}

Conn:
	conn, err := database.Open()
	if err != nil {
		log.WithError(err).Error("Couldn't establish DB connection. Retrying in 3 seconds...")

		// Retry connection until succeeds
		time.Sleep(3 * time.Second)
		goto Conn
	}

	log.Info("Migrating account service models")
	// Do migrations
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	conn.AutoMigrate(
		&Account{},
		&PasswordHash{},
		&Profile{},
		&Qualification{},
		&WorkExperience{},
		&Lesson{},
		&ResourceMetadata{},
		&ResourceData{},
		&Subject{},
		&SubjectTaught{},
	)
	// Add some test users so we don't need to manually test things
	CreateDebugData()

}
