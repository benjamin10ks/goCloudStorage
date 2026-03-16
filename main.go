package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

type App struct {
	db      *sql.DB
	tmpl    *template.Template
	auth    *AuthService
	storage *StorageService
}

func main() {
	db, err := sql.Open("sqlite", "./cloud.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Println("Database connection established")

	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	app := &App{
		db:      db,
		tmpl:    template.Must(template.ParseGlob("web/templates/*.tmpl")),
		auth:    NewAuthService(),
		storage: NewStorageService(BASE_PATH),
	}

	mux := http.NewServeMux()
	server := newServer(mux)

	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	// Dashboard routes
	mux.HandleFunc("/", app.handleLandingPage)
	mux.HandleFunc("/home", app.handleHomePage)

	// Auth routes
	mux.HandleFunc("/auth/register", app.handleRegister)
	mux.HandleFunc("/auth/login", app.handleLogin)
	mux.HandleFunc("/auth/logout", app.handleLogout)

	// User routes
	mux.HandleFunc("GET /api/user/{id}", app.handleUserProfile)
	mux.HandleFunc("PUT /api/user/{id}", app.handleUpdateUser)
	mux.HandleFunc("DELETE /api/user/{id}", app.handleDeleteUser)
	mux.HandleFunc("GET /api/user/{id}/files", app.handleListUserFiles)

	// File routes
	mux.HandleFunc("POST /api/upload", app.handleUpload)
	mux.HandleFunc("GET /api/files/{id}/download", app.handleDownload)
	mux.HandleFunc("DELETE /api/delete/{id}", app.handleDeleteFile)
	mux.HandleFunc("POST /api/files/{id}/share", app.handleShareFile)

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
