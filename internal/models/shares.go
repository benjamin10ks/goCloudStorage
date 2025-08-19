package models

type Shares struct {
	FileID     string `json:"file_id"`
	ShareBy    string `json:"share_id"`
	SharedWith string `json:"shared_with"`
	Permission string `json:"permission"`
}
