package services

import (
	"database/sql"

	"github.com/benjamin10ks/goCloudStorage/internal/models"
	"golang.org/x/crypto/bcrypt"
)

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

func (u *UserService) CreateUser(username, password string) (*models.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO users (user_name, password_hash) VALUES ($1, $2) RETURNING user_id`
	userID := 0
	u.db.QueryRow(query, username, passwordHash).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:           userID,
		Name:         username,
		PasswordHash: passwordHash,
	}, nil
}

func (u *UserService) DeleteUser(userID int) error {
	// Placeholder implementation
	query := `DELETE FROM users WHERE user_id = $1`
	return nil
}

func hashPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return passwordHash, nil
}
