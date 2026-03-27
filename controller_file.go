package main

import (
	"log"
	"net/http"
)

// TODO: implement file upload, download, delete, and sharing handlers
func (a *App) handleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from form: %v", err)
	}
	log.Printf("Received file upload: %s (%d bytes)", header.Filename, header.Size)
	defer file.Close()

	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		log.Printf("File upload failed: %v", err)
		return
	}
}

func (a *App) handleDownload(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleShareFile(w http.ResponseWriter, r *http.Request) {
}
