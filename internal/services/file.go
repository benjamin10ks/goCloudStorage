// Package services contains management services for users, files, and directories.
package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
)

type FileService struct {
	DB      *sql.DB
	storage *StorageService
}

func NewFileService(db *sql.DB, storage *StorageService) *FileService {
	return &FileService{
		db,
		storage,
	}
}

func detectType(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)

	return http.DetectContentType(buf[:n]), nil
}

func (fs *FileService) ApplyFileType(filepath string) (string, error) {
	fileType, err := detectType(filepath)
	if err != nil {
		fmt.Printf("failed to detect file type, %d", err)
	}

	return fileType, nil
}
