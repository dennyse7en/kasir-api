package database

import (
	"database/sql"
	"fmt"
	"log"
)

func NewPostgresDB(driver, source string) (*sql.DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}
