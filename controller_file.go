package main

import (
	"log"
	"net/http"
)

func (a *App) handleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from form: %v", err)
	}
	defer file.Close()

	path, err := a.storage.UploadFile(userID, file, header)
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
