package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
)

type DB struct {
	*gorm.DB
}

// NewDB returns a new gorm instance
func NewDB(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(path, "gorm.db")), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&User{},
		&CustomVariable{},
		&Stat{},
		&FallbackQueue{},
		&CustomCommand{},
		&Client{},
	)
	if err != nil {
		return nil, err
	}
	return &DB{db}, err
}
