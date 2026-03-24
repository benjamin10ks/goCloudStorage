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
		auth:    NewAuthService(db),
		storage: NewStorageService(BASE_PATH, db),
	}

	mux := http.NewServeMux()
	server := newServer(mux)

	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	// Dashboard routes
	mux.HandleFunc("GET /", app.handleLandingPage)
	mux.HandleFunc("GET /home", app.handleHomePage)

	// display login and registration pages
	mux.HandleFunc("GET /auth/register", app.handleRegister)
	mux.HandleFunc("GET /auth/login", app.handleLogin)
	mux.HandleFunc("POST /auth/logout", app.handleLogout)

	// primary auth routes for passkey and github oauth
	mux.HandleFunc("POST /auth/passkey/register/begin", app.handlePasskeyBeginRegister)
	mux.HandleFunc("POST /auth/passkey/register/finish", app.handlePasskeyFinishRegister)
	mux.HandleFunc("POST /auth/passkey/login/begin", app.handlePasskeyBeginLogin)
	mux.HandleFunc("POST /auth/passkey/login/finish", app.handlePasskeyFinishLogin)

	mux.HandleFunc("POST /auth/github", app.handleGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", app.handleGitHubCallback)

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
