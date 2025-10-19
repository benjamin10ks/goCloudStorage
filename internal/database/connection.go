// Package database manages the database connection and provides functions to interact with the database.
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
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	connStr := fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=disable", host, user, port, password, dbname)

	log.Printf("Connecting to database with user: %s, dbname: %s", user, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	maxRetries := 5
	for i := range maxRetries {
		err := db.Ping()
		if err == nil {
			log.Println("Successfully connected to the database")
			return &PostgresStore{DB: db}, nil
		}
		log.Printf("Failed to connect to database retrying... (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
		db.SetConnMaxLifetime(5 * time.Minute)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{DB: db}, nil
}
