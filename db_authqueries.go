package main

import (
	"context"
	"database/sql"
	"fmt"
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

func updatePasskeySignCount(ctx context.Context, db *sql.DB, userID int64, credentialID []byte, newCount uint32) error {
	if newCount == 0 {
		return nil
	}

	var currentCount uint32
	err := db.QueryRowContext(
		ctx, "SELECT sign_count FROM passkeys WHERE credential_id = ?", credentialID,
	).Scan(&currentCount)
	if err != nil {
		return err
	}

	if newCount <= currentCount {
		return fmt.Errorf("sign count invalid: possible credential clone detected")
	}

	_, err = db.ExecContext(
		ctx, "UPDATE passkeys SET sign_count = ?, last_used_at = CURRENT_TIMESTAMP WHERE credential_id = ?",
		newCount, credentialID,
	)
	return err
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

func deleteSession(ctx context.Context, db *sql.DB, sessionToken string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", sessionToken)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func storeGithubToken(ctx context.Context, db *sql.DB, userID int64, githubUserID, token string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO oauth_accounts (user_id, provider, provider_user_id, access_token)
		VALUES (?, 'github', ?, ?)
		ON CONFLICT (provider, provider_user_id)
		DO UPDATE SET access_token = excluded.access_token`,
		userID, userID, token,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
