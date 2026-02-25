package database

import (
	"sika/internal/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

func NewConnection(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	instance := &Database{
		Db: db,
	}

	return instance, nil
}

func (db *Database) Migrate() error {
	return db.Db.AutoMigrate(&models.User{}, &models.Address{})
}
