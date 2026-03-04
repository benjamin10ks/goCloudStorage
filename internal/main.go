package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("/api/upload", handleUpload)
	mux.HandleFunc("/api/download", handleDownload)
	mux.HandleFunc("/api/list", handleListFiles)
	mux.HandleFunc("/api/delete", handleDeleteFile)
	mux.HandleFunc("/api/share", handleShareFile)

	server := newServer(mux)

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
