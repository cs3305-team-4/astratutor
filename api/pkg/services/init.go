package services

import (
	"fmt"

	"github.com/cs3305-team-4/api/pkg/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var conn *gorm.DB

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Could not read config file %s\n", err))
	}
	conn, err = db.DB()
	if err != nil {
		log.WithError(err).Error("Couldn't establish DB connection")
	}
	modelMigrations()
}

func modelMigrations() {
	log.Info("Migrating account service models")
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	conn.AutoMigrate(&Account{}, &PasswordHash{}, &Profile{}, &Qualification{}, &WorkExperience{})
}
