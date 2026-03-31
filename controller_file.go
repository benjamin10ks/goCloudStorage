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

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close uploaded file: %v", err)
			return
		}
	}()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("No session cookie found: %v", err)
		http.Error(w, "Unauthorized: no session cookie", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	user, err := getUserBySessionToken(ctx, a.db, cookie.Value)
	if err != nil {
		log.Printf("Failed to get user from session: %v", err)
		http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
		return
	}
	log.Printf("Authenticated user: %s (ID: %d)", user.Username, user.ID)

	path, err := a.storage.UploadFile(user.ID, file, header)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		log.Printf("File upload failed: %v", err)
		return
	}
	log.Printf("File uploaded successfully: %s", path)
}

func (a *App) handleDownload(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleShareFile(w http.ResponseWriter, r *http.Request) {
}
