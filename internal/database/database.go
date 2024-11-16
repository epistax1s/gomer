package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("/app/database/gomer.db"), &gorm.Config{})
}
