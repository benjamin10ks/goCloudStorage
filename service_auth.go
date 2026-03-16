package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Register(token string, db *sql.DB) (User, error) {
	return User{}, nil
}

func (a *AuthService) Login(token string) error {
	return nil
}

func (a *AuthService) Logout() error {
	return nil
}

func generateToken() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
