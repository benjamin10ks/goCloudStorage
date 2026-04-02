package main

import (
	"context"
	"database/sql"
)

func savePathToDB(ctx context.Context, db *sql.DB, userID int64, filename string, path string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO files (user_id, filename, path) VALUES (?, ?, ?)", userID, filename, path)
	if err != nil {
		return err
	}
	return tx.Commit()
}
