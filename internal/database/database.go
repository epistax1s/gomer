package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("./database/gomer.db"), &gorm.Config{})
}
