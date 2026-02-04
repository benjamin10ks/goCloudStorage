// Package handlers holds all main routes handlers
package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")

	dst, err := os.Create()
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer func() {
		if cerr := dst.Close(); cerr != nil {
			err = cerr
		}
	}()

	_, err = io.Copy(dst, r.Body)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received file with Content-Type: %s\n", ct)

	w.WriteHeader(http.StatusOK)
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Placeholder for download handler logic
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Download handler not implemented yet."))
}

func HandleListFiles(w http.ResponseWriter, r *http.Request) {
	// Placeholder for list files handler logic
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "List files handler not implemented yet."}`))
}

func HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	// Placeholder for delete file handler logic
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete file handler not implemented yet."))
}

func HandleShareFile(w http.ResponseWriter, r *http.Request) {
	// Placeholder for delete file handler logic
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete file handler not implemented yet."))
}
