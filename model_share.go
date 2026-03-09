package main

type Share struct {
	ID         int `json:"id"`
	FileID     int `json:"file_id"`
	ToUserID   int `json:"to_user_id"`
	FromUserID int `json:"from_user_id"`
}
