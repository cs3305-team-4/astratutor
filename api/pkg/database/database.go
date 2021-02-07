package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var db *gorm.DB

func Open() (*gorm.DB, error) {
	var err error

	if db == nil {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.host"),
			viper.GetString("database.port"),
			viper.GetString("database.database"),
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	return db, err
}
