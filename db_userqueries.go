package main

import (
	"context"
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

func getUserByCredentialID(db *sql.DB, credentialID []byte) (*User, error) {
	var user User
	err := db.QueryRow("SELECT u.id, u.username, u.display_name FROM users u JOIN passkeys p ON u.id = p.user_id WHERE p.credential_id = ?", credentialID).Scan(&user.ID, &user.Username, &user.DisplayName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserByGithubID(ctx context.Context, db *sql.DB, githubUserID string) (*User, error) {
	row := db.QueryRowContext(ctx, `
		SELECT u.id, u.username
		FROM users u
		JOIN oauth_accounts oa ON oa.user_id = u.id
		WHERE oa.provider = 'github' AND oa.provider_user_id = ?`,
		githubUserID,
	)
	var u User
	if err := row.Scan(&u.ID, &u.Username); err != nil {
		return nil, err
	}
	return &u, nil
}

func createUserFromGithub(ctx context.Context, db *sql.DB, githubUserID, accessToken string) (*User, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `INSERT INTO users (username) VALUES (?)`, "github_"+githubUserID)
	if err != nil {
		return nil, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO oauth_accounts (user_id, provider, provider_user_id, access_token)
		VALUES (?, 'github', ?, ?)`,
		userID, githubUserID, accessToken,
	)
	if err != nil {
		return nil, err
	}

	return &User{ID: userID, Username: "github_" + githubUserID}, tx.Commit()
}

func updateUserDisplayName(ctx context.Context, db *sql.DB, userID int64, displayName string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE users SET display_name = ? WHERE id = ?`, displayName, userID)
	if err != nil {
		return err
	}
	return nil
}
