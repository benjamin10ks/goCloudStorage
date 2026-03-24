package main

import (
	"database/sql"
)

func createUser(db *sql.DB, username string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	var userID int64
	err = tx.QueryRow("INSERT INTO users (username) VALUES (?) RETURNING id", username).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return userID, nil
}

func getUserByID(db *sql.DB, userID string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT username, display_name FROM users WHERE id = ?", userID).Scan(&user.Username, &user.DisplayName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(db *sql.DB, username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, display_name FROM users WHERE username = ?", username).Scan(&user.ID, &user.DisplayName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
