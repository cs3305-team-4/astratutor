package services

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cs3305-team-4/api/pkg/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/validator.v2"
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

// SetCustomValidator will set field validators.
func SetCustomValidators() {
	validator.SetValidationFunc("passwd", func(v interface{}, param string) error {
		st := reflect.ValueOf(v)
		if st.Kind() != reflect.String {
			return errors.New("passwd only validates strings")
		}
		val := st.String()
		if len(val) < 8 {
			return errors.New("Must have at least 8 characters")
		}
		if strings.ToLower(val) == val {
			return errors.New("Must have at least one upper case letter")
		}
		if strings.ToUpper(val) == val {
			return errors.New("Must have at least one lower case letter")
		}
		if strings.ToUpper(val) == val {
			return errors.New("Must have at least one lower case letter")
		}
		return nil
	})
}
