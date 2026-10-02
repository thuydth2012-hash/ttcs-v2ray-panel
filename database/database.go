package database

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DefaultDBPath = "/etc/ttcs/ttcs.db"

var DB *gorm.DB

func Init(dbPath string) error {
	if dbPath == "" {
		dbPath = DefaultDBPath
	}

	dir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("cannot open database: %w", err)
	}

	DB = db

	return nil
}
