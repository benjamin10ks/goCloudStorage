package main

import (
	"context"
	"database/sql"
	"log"
)

func savePathToDB(ctx context.Context, db *sql.DB, userID int64, filename string, path string) (int64, error) {
	return saveFileWithSize(ctx, db, userID, filename, path, 0)
}

func saveFileWithSize(ctx context.Context, db *sql.DB, userID int64, filename string, path string, size int64) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	result, err := tx.ExecContext(ctx, "INSERT INTO files (user_id, filename, filepath, size, type) VALUES (?, ?, ?, ?, '')", userID, filename, path, size)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	fileID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return fileID, nil
}

func getUserFiles(ctx context.Context, db *sql.DB, userID int64) ([]File, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT id, filename, filepath, type, size, created_at FROM files WHERE user_id = ? ORDER BY created_at DESC", userID)
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
		if err := rows.Scan(&file.ID, &file.FileName, &file.FilePath, &file.Type, &file.Size, &file.FileMetadata.CreatedAt); err != nil {
			log.Printf("Error scanning file row: %v", err)
			continue
		}
		files = append(files, file)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return files, nil
}

func getFilepathByID(ctx context.Context, db *sql.DB, userID int64, fileID string) (string, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var filepath string
	err = tx.QueryRowContext(ctx, "SELECT filepath FROM files WHERE id = ? AND user_id = ?", fileID, userID).Scan(&filepath)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return filepath, nil
}

func deleteFileByID(ctx context.Context, db *sql.DB, userID int64, fileID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM files WHERE id = ? AND user_id = ?", fileID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func deleteFileByPath(ctx context.Context, db *sql.DB, path string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM files WHERE filepath = ?", path)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
