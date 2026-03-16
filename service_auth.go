package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthService) Register(email, password string, db *sql.DB) (User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return User{}, err
	}

	var user User
	err = db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, hashedPassword).Scan(&user.ID)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return User{}, err
	}
	return user, nil
}

func (a *AuthService) Login(email, password string) error {
	passwordHash := ""
	err := db.QueryRow("SELECT password_hash FROM users WHERE email = $1", email).Scan(&passwordHash)
	err = verifyPassword(password, passwordHash)
	return nil
}

func (a *AuthService) Logout() error {
	return nil
}
