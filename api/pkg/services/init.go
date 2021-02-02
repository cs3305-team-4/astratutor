package services

import (
	"github.com/cs3305-team-4/api/pkg/db"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var conn *gorm.DB

func init() {
	var err error
	conn, err = db.DB()
	if err != nil {
		log.WithError(err).Error("Couldn't establish DB connection")
	}
}
