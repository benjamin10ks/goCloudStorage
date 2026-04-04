package main

import (
	"context"
	"database/sql"
	"log"
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

func getUserFiles(ctx context.Context, db *sql.DB, userID int64) ([]File, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, "SELECT filename, type, size, updated_at FROM files WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	var files []File
	for rows.Next() {
		var file File
		if err := rows.Scan(&file.FileName, &file.Type, &file.Size, &file.FileMetadata); err != nil {
			log.Printf("Error scanning file row: %v", err)
			continue
		}
		files = append(files, file)
	}
	return files, nil
}
