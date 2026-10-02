package database

import "fmt"

func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	return DB.AutoMigrate(
		&User{},
		&Setting{},
		&Inbound{},
	)
}
