// package database manages the database connection and provides functions to interact with the database.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbname, user, password)

	log.Printf("Connecting to database with user: %s, dbname: %s", user, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresStore{db: db}, nil
}
