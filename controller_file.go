package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (a *App) handleUploadFile(w http.ResponseWriter, r *http.Request) {
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

	fileType := header.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}
	log.Printf("Determined file type: %s", fileType)

	fileID, err := saveFileWithSize(ctx, a.db, user.ID, header.Filename, path, header.Size, fileType)
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
		ID:       fileID,
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

func (a *App) handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		log.Printf("User not found in context")
		http.Error(w, "Unauthorized: user not found in context", http.StatusUnauthorized)
		return
	}

	fileID := r.PathValue("id")

	ctx := r.Context()

	file, err := getFileByID(ctx, a.db, user.ID, fileID)
	if err != nil {
		log.Printf("Failed to retrieve file: %v", err)
		http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
		return
	}
	log.Printf("Retrieved file for download: %s (ID: %d)", file.FileName, file.ID)

	f, err := os.Open(file.FilePath)
	if err != nil {
		log.Printf("Failed to open file for download: %v", err)
		http.Error(w, "Failed to open file for download", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	log.Printf("Opened file for download: %s", file.FilePath)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", file.FileName))
	w.Header().Set("Content-Type", file.Type)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
	log.Printf("Set response headers for file download: Content-Disposition=attachment; filename=%q, Content-Type=%s, Content-Length=%d", file.FileName, file.Type, file.Size)

	if _, err := io.Copy(w, f); err != nil {
		log.Printf("Failed to send file to client: %v", err)
	}
}

func (a *App) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		log.Printf("User not found in context")
		http.Error(w, "Unauthorized: user not found in context", http.StatusUnauthorized)
		return
	}

	fileID := r.PathValue("id")

	ctx := r.Context()

	filepath, err := getFilepathByID(ctx, a.db, user.ID, fileID)
	if err != nil {
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Failed to retrieve file path for deletion", "type": "error"}}`)
		log.Printf("Failed to retrieve file path for deletion: %v", err)
		http.Error(w, "Failed to retrieve file path for deletion", http.StatusInternalServerError)
		return
	}
	log.Printf("Retrieved file path for deletion: %s", filepath)

	err = deleteFileByID(ctx, a.db, user.ID, fileID)
	if err != nil {
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Failed to delete file metadata", "type": "error"}}`)
		log.Printf("Failed to delete file metadata: %v", err)
		http.Error(w, "Failed to delete file metadata", http.StatusInternalServerError)
		return
	}
	log.Printf("Deleted file metadata for file ID %s", fileID)

	err = a.storage.DeleteFile(filepath)
	if err != nil {
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Failed to delete file", "type": "error"}}`)
		log.Printf("Failed to delete file: %v", err)
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}
	log.Printf("Deleted file from storage: %s", filepath)
}

func (a *App) handleShareFile(w http.ResponseWriter, r *http.Request) {
}
