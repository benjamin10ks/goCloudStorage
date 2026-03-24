package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type AuthService struct {
	webauthn *webauthn.WebAuthn
	db       *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	wauthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "My Cloud Storage",
		RPID:          "localhost",                       // Adjust when deploying
		RPOrigins:     []string{"http://localhost:8080"}, // Adjust when deploying
	})
	log.Fatalf("Failed to initialize WebAuthn: %v", err)
	return &AuthService{webauthn: wauthn, db: db}
}

func (a *AuthService) BeginRegistration(user *User) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	options, sessionData, err := a.webauthn.BeginRegistration(user)
	if err != nil {
		return nil, nil, err
	}
	return options, sessionData, nil
}

func (a *AuthService) FinishRegistration(user *User, sessionData *webauthn.SessionData, r *http.Request) error {
	credential, err := a.webauthn.FinishRegistration(user, *sessionData, r)
	if err != nil {
		return err
	}
	return savePasskey(a.db, user.ID, credential)
}

func (a *AuthService) BeginLogin(user *User) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	options, sessionData, err := a.webauthn.BeginLogin(user)
	if err != nil {
		return nil, nil, err
	}
	return options, sessionData, nil
}

func (a *AuthService) FinishLogin(user *User, sessionData *webauthn.SessionData, r *http.Request) error {
	credential, err := a.webauthn.FinishLogin(user, *sessionData, r)
	if err != nil {
		return err
	}
	return updatePasskeySignCount(a.db, user.ID, credential.ID, credential.Authenticator.SignCount)
}

func generateToken() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
