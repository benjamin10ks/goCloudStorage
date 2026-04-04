package main

import (
	"log"
	"net/http"
	"time"
)

func (a *App) handleUpload(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		log.Printf("User not found in context")
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Unauthorized", "type": "error"}}`)
		http.Error(w, "Unauthorized: user not found in context", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from form: %v", err)
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Failed to read uploaded file", "type": "error"}}`)
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	log.Printf("Received file upload: %s (%d bytes)", header.Filename, header.Size)

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close uploaded file: %v", err)
			return
		}
	}()

	path, err := a.storage.UploadMultipartFile(user.ID, file, header)
	if err != nil {
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Failed to upload file", "type": "error"}}`)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		log.Printf("File upload failed: %v", err)
		return
	}
	log.Printf("File uploaded successfully: %s", path)

	ctx := r.Context()

	fileID, err := saveFileWithSize(ctx, a.db, user.ID, header.Filename, path, header.Size)
	if err != nil {
		log.Printf("Failed to save file path to DB: %v", err)
		err = a.storage.DeleteFile(path)
		if err != nil {
			log.Printf("Failed to delete file after DB save failure: %v", err)
			w.Header().Set("HX-Trigger", `{"showToast": {"message": "Failed to save file metadata and cleanup file", "type": "error"}}`)
			http.Error(w, "Failed to save file metadata and cleanup file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Failed to save file metadata", "type": "error"}}`)
		http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
		return
	}
	log.Printf("File metadata saved to DB for user %d: %s", user.ID, header.Filename)

	w.Header().Set("HX-Trigger", `{"showToast": {"message": "File uploaded successfully!", "type": "success"}, "fileUploaded": true}`)
	w.WriteHeader(http.StatusCreated)

	uploadedFile := &File{
		ID:       int(fileID),
		FileName: header.Filename,
		FilePath: path,
		Size:     header.Size,
		Type:     "",
		FileMetadata: FileMetadata{
			CreatedAt: time.Now(),
		},
	}

	err = a.tmpl["dashboard"].ExecuteTemplate(w, "file_card", uploadedFile)
	if err != nil {
		log.Printf("Error rendering file card: %v", err)
	}
}

func (a *App) handleDownload(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleShareFile(w http.ResponseWriter, r *http.Request) {
}
