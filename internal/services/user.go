package services

import "database/sql"

type UserService struct {
	db      *sql.DB
	storage *StorageService
}

func NewUserService(db *sql.DB, storage *StorageService) *UserService {
	return &UserService{
		db:      db,
		storage: storage,
	}
}

func (u *UserService) CreateUser(username, passwordHash string) (int, error) {
	// Placeholder implementation
	return 0, nil
}

func (u *UserService) DeleteUser(userID int) error {
	// Placeholder implementation
	return nil
}
