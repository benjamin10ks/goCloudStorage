package main

import (
	"log"
	"net/http"

	"github.com/benjamin10ks/goCloudStorage/internal/config"
	"github.com/benjamin10ks/goCloudStorage/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.PlaceHolderHandler)

	server := config.NewServer(mux)

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
