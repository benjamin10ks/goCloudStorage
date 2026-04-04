package main

import (
	"fmt"
	"time"
)

type File struct {
	ID           int64        `json:"id"`
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

func (f *File) Name() string {
	return f.FileName
}

func (f *File) SizeFormatted() string {
	const unit = 1024
	if f.Size < unit {
		return fmt.Sprintf("%d B", f.Size)
	}
	div, exp := int64(unit), 0
	for n := f.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(f.Size)/float64(div), "KMGTPE"[exp])
}

func (f *File) CreatedAt() string {
	return f.FileMetadata.CreatedAt.Format("Jan 2, 2006")
}
