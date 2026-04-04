package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StorageService struct {
	basePath string
}

func NewStorageService(basePath string, db *sql.DB) *StorageService {
	return &StorageService{
		basePath: basePath,
	}
}

func (s *StorageService) UploadFile(userID int64, r io.Reader, filename string, size int64) (string, error) {
	filename = filepath.Base(filepath.Clean(filename))
	if filename == "." || filename == "/" {
		return "", fmt.Errorf("invalid filename: %s", filename)
	}

	userDir := filepath.Join(s.basePath, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	destPath := filepath.Join(userDir, filename)
	destPath = deduplicatePath(destPath)

	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}

	defer func() {
		closeErr := dst.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("failed to close destination file: %w", closeErr)
		}
	}()

	written, err := io.Copy(dst, r)
	if err != nil {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if size > 0 && written != size {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("written size mismatch: expected %d bytes, got %d bytes", size, written)
	}

	return destPath, nil
}

func (s *StorageService) UploadMultipartFile(userID int64, file multipart.File, header *multipart.FileHeader) (string, error) {
	return s.UploadFile(userID, file, header.Filename, header.Size)
}

func (s *StorageService) DeleteFile(path string) error {
	return nil
}

func (s *StorageService) getFileType(file *multipart.File) string {
	fileBytes, err := io.ReadAll(*file)
	if err != nil {
		log.Printf("Failed to read file for type detection: %v", err)
	}
	return http.DetectContentType(fileBytes)
}

func deduplicatePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	ext := filepath.Ext(path)
	nameWithoutExt := strings.TrimSuffix(path, ext)

	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s (%d)%s", nameWithoutExt, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
