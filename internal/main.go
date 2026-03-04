package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

type App struct {
	db   *sql.DB
	tmpl *template.Template
}

func main() {
	db, err := sql.Open("sqlite", "./cloud.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Println("Database connection established")

	app := &App{
		db:   db,
		tmpl: template.Must(template.ParseGlob("templates/*.tmpl")),
	}

	mux := http.NewServeMux()
	server := newServer(mux)

	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	mux.HandleFunc("/", handleLandingPage)

	// Auth routes
	mux.HandleFunc("/auth/register", handleRegister)
	mux.HandleFunc("/auth/login", handleLogin)
	mux.HandleFunc("/auth/logout", handleLogout)

	// User routes
	mux.HandleFunc("GET /api/user/{id}", handleUserProfile)
	mux.HandleFunc("PUT /api/user/{id}", handleUpdateUser)
	mux.HandleFunc("DELETE /api/user/{id}", handleDeleteUser)
	mux.HandleFunc("GET /api/user/{id}/files", handleListUserFiles)

	// File routes
	mux.HandleFunc("POST /api/upload", handleUpload)
	mux.HandleFunc("GET /api/files/{id}/download", handleDownload)
	mux.HandleFunc("DELETE /api/delete/{id}", handleDeleteFile)
	mux.HandleFunc("POST /api/files/{id}/share", handleShareFile)

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
