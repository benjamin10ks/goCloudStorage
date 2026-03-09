package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type StorageService struct {
	basePath string
}

func NewStorageService(basePath string) *StorageService {
	return &StorageService{
		basePath: basePath,
	}
}

func (s *StorageService) UploadFile(userID int, file multipart.File, header *multipart.FileHeader) (string, error) {
	return "", nil
}

func (s *StorageService) getFileType(file *multipart.File) string {
	fileBytes, err := io.ReadAll(*file)
	if err != nil {
		log.Printf("Failed to read file for type detection: %v", err)
	}
	return http.DetectContentType(fileBytes)
}
