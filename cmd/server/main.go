package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/benjamin10ks/goCloudStorage/internal/config"
	"github.com/benjamin10ks/goCloudStorage/internal/database"
	"github.com/benjamin10ks/goCloudStorage/internal/handlers"
	"github.com/benjamin10ks/goCloudStorage/internal/services"
)

func main() {
	store, err := database.NewPostgresStore()
	storageService := services.NewStorageService("/var/goCloudStorage/storage/users") // temp path
	userService := services.NewUserService(store.DB, storageService)

	h := handlers.NewHandler(userService)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Printf("Connected to database successfully %+v\n", store)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	mux.HandleFunc("/", handlers.TemplateRenderHandler)

	mux.HandleFunc("/api/upload", handlers.HandleUpload)
	mux.HandleFunc("/api/download", handlers.HandleDownload)
	mux.HandleFunc("/api/list", handlers.HandleListFiles)
	mux.HandleFunc("/api/delete", handlers.HandleDeleteFile)
	mux.HandleFunc("/api/share", handlers.HandleShareFile)

	server := config.NewServer(mux)

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
