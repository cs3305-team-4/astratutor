package services

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/cs3305-team-4/api/pkg/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/validator.v2"
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

	conn, err := database.Open()
	if err != nil {
		log.WithError(err).Error("Couldn't establish DB connection")
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
	)

	// Add some test users so we don't need to manually test things
	CreateDebugData()
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
		numRe := regexp.MustCompile(`[0-9]+`)
		if !numRe.MatchString(val) {
			return errors.New("Must have at least one number")
		}
		return nil
	})
}
