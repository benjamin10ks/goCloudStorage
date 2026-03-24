package main

import (
	"database/sql"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

// this prob wont work as is
func getWebAuthnSession(db *sql.DB, reqType string, userID int64) (*webauthn.SessionData, error) {
	var sessionData webauthn.SessionData
	result := db.QueryRow("SELECT session_data FROM webauthn_sessions WHERE user_id = ? AND auth_method = ? AND expires_at > ?", userID, reqType, time.Now().Unix())
	result.Scan(&sessionData)
	return &sessionData, nil
}

func saveWebAuthnSession(db *sql.DB, reqType string, userID int64, sessionData *webauthn.SessionData) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(60 * time.Second).Unix()
	_, err = tx.Exec("INSERT INTO webauthn_sessions (user_id, auth_method, expires_at) VALUES (?, ?, ? ,?)", userID, reqType, sessionData, expiresAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// temporary function to save passkey credential
func savePasskey(db *sql.DB, userID int64, credential *webauthn.Credential) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO passkeys (user_id, credential_id, public_key, sign_count) VALUES (?, ?, ?, ?)", userID, credential.ID, credential.PublicKey, credential.Authenticator.SignCount)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func updatePasskeySignCount(db *sql.DB, userID int64, credentialID []byte, signCount uint32) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE passkeys SET sign_count = ? WHERE user_id = ? AND credential_id = ?", signCount, userID, credentialID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func saveSession(db *sql.DB, userID int64, sessionToken string, authMethod string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	_, err = tx.Exec("INSERT INTO sessions (id, user_id, auth_method,expires_at) VALUES (?, ?, ? ,?)", sessionToken, userID, authMethod, expiresAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
