package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetGormConnection creates gorm database śtruct
func GetGormConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	return db, err
}
