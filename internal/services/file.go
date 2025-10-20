// Package services contains management services for users, files, and directories.
package services

import "database/sql"

type FileService struct {
	DB      *sql.DB
	storage *StorageService
}
