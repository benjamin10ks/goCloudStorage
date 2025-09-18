package services

import (
	"fmt"
	"os"
	"path/filepath"
)

type StorageService struct {
	basePath string
}

func NewStorageService(basePath string) *StorageService {
	return &StorageService{
		basePath: basePath,
	}
}

func (s *StorageService) CreateUserDirectory(userID int) error {
	err := os.MkdirAll(s.GetUserDirectory(userID), 0700)
	if err != nil {
		return fmt.Errorf("failed to create user directory: %w", err)
	}
	return nil
}

func (s *StorageService) DeleteUserDirectory(userID int) error {
	err := os.RemoveAll(s.GetUserDirectory(userID))
	if err != nil {
		return fmt.Errorf("failed to delete user directory: %w", err)
	}
	return nil
}

func (s *StorageService) GetUserDirectory(userID int) string {
	return filepath.Join(s.basePath, "users", fmt.Sprintf("%d", userID))
}

func (s *StorageService) GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	name := originalName[0 : len(originalName)-len(ext)]
	uniqueName := fmt.Sprintf("%s_%d%s", name, os.Getpid(), ext)
	return uniqueName
}

func (s *StorageService) EnsureUserDirExsits(userID int) error {
	userDir := s.GetUserDirectory(userID)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		return s.CreateUserDirectory(userID)
	}
	return nil
}
