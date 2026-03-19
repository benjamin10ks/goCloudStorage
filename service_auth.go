package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"

	"github.com/go-webauthn/webauthn/webauthn"
)

type AuthService struct {
	webauthn *webauthn.WebAuthn
	db       *sql.DB
}

func NewAuthService(db *sql.DB) (*AuthService, error) {
	wauthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "My Cloud Storage",
		RPID:          "localhost",                       // Adjust when deploying
		RPOrigins:     []string{"http://localhost:8080"}, // Adjust when deploying
	})
	if err != nil {
		return nil, err
	}
	return &AuthService{webauthn: wauthn, db: db}, nil
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
