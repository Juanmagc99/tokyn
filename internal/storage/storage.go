package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(name string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
