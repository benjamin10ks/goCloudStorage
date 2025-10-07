package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/benjamin10ks/goCloudStorage/internal/config"
	"github.com/benjamin10ks/goCloudStorage/internal/database"
	"github.com/benjamin10ks/goCloudStorage/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	store, err := database.NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Printf("Connected to database successfully %+v\n", store)

	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	mux.HandleFunc("/", handlers.TemplateRenderHandler)

	mux.HandleFunc("/api/upload", handlers.UploadHandler)
	mux.HandleFunc("/api/download", handlers.DownloadHandler)
	mux.HandleFunc("/api/list", handlers.ListFilesHandler)
	mux.HandleFunc("/api/delete", handlers.DeleteFileHandler)
	mux.HandleFunc("/api/share", handlers.ShareFileHandler)

	server := config.NewServer(mux)

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
