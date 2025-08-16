// Package models contains User and File structs
package models

type User struct {
	ID   int64  `json:"id" db:"user_id"`
	Name string `json:"name" db:"user_name"`
}
