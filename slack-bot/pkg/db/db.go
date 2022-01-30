package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
)

type DB struct {
	*gorm.DB
}

func NewDB(path string) (*DB, error) {
	// github.com/mattn/go-sqlite3
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
	)
	if err != nil {
		return nil, err
	}
	return &DB{db}, err
}
