// Package models contains User and File structs
package models

type User struct {
	ID           int    `json:"id" db:"user_id"`
	Name         string `json:"name" db:"user_name"`
	PasswordHash []byte
}

type CreateAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
