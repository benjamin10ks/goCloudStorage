package main

import "time"

type File struct {
	ID           int          `json:"id"`
	FileName     string       `json:"filename"`
	FilePath     string       `json:"filepath"`
	Type         string       `json:"type"`
	Size         int64        `json:"size"`
	FileMetadata FileMetadata `json:"metadata"`
}

type FileMetadata struct {
	CreatedAt  time.Time `json:"created_at"`
	AccessedAt time.Time `json:"accessed_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
