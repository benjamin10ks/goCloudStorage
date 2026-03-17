package main

type Account struct {
	ID int `json:"id"`
}

type User struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
	BasePath    string `json:"base_path"`
}
