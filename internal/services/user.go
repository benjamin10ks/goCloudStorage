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
	sql := `INSERT INTO users (user_name, password_hash) VALUES ($1, $2) RETURNING user_id`
	return 0, nil
}

func (u *UserService) DeleteUser(userID int) error {
	// Placeholder implementation
	sql := `DELETE FROM users WHERE user_id = $1`
	return nil
}
