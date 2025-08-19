package models

type File struct {
	ID       string `json:"id" db:"file_id"`
	UserID   int64  `json:"user_id" db:"user_id"`
	Filename string `json:"filename" db:"filename"`
	FilePath string `json:"file_path" db:"file_path"`
}
