package services

import (
	"database/sql"
	"log"

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
	err = u.db.QueryRow(query, username, passwordHash).Scan(&userID)
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
	query := `DELETE FROM users WHERE user_id = $2`
	res, err := u.db.Exec(query, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}
	return nil
}

func hashPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return passwordHash, nil
}
