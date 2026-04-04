package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

type App struct {
	db      *sql.DB
	tmpl    map[string]*template.Template
	auth    *AuthService
	storage *StorageService
}

func main() {
	if GithubClientID == "" || GithubSecret == "" {
		log.Fatal("GitHub OAuth credentials are not set. Please set GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET environment variables.")
	}

	db, err := sql.Open("sqlite", "./cloud.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Println("Database connection established")

	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	templates, err := loadTemplates()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	app := &App{
		db:      db,
		tmpl:    templates,
		auth:    NewAuthService(db),
		storage: NewStorageService(BASE_PATH, db),
	}

	app.db.SetMaxOpenConns(1)
	mux := http.NewServeMux()
	server := newServer(mux)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	// Dashboard routes
	mux.HandleFunc("GET /", app.handleLandingPage)
	mux.HandleFunc("GET /home", app.requireAuth(app.handleHomePage))

	// display login and registration pages
	mux.HandleFunc("GET /auth/register", app.handleRegister)
	mux.HandleFunc("GET /auth/login", app.handleLogin)
	mux.HandleFunc("POST /auth/logout", app.handleLogout)

	// primary auth routes for passkey and github oauth
	mux.HandleFunc("POST /auth/register/username", app.handleRegisterUsername)
	mux.HandleFunc("POST /auth/register/passkey/begin/{id}", app.handlePasskeyBeginRegister)
	mux.HandleFunc("POST /auth/register/passkey/finish/{id}", app.handlePasskeyFinishRegister)
	mux.HandleFunc("POST /auth/login/passkey/begin/{id}", app.handlePasskeyBeginLogin)
	mux.HandleFunc("POST /auth/login/passkey/finish/{id}", app.handlePasskeyFinishLogin)

	mux.HandleFunc("POST /auth/github", app.handleGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", app.handleGitHubCallback)

	// User routes
	mux.HandleFunc("GET /api/user/{id}", app.handleUserProfile)
	mux.HandleFunc("PUT /api/user/{id}", app.requireAuth(app.handleUpdateUser))
	mux.HandleFunc("DELETE /api/user/{id}", app.requireAuth(app.handleDeleteUser))
	mux.HandleFunc("GET /api/user/{id}/files", app.requireAuth(app.handleListUserFiles))

	// File routes
	mux.HandleFunc("POST /api/files/upload", app.requireAuth(app.handleUploadFile))
	mux.HandleFunc("GET /api/files/{id}/download", app.requireAuth(app.handleDownloadFile))
	mux.HandleFunc("DELETE /api/files/delete/{id}", app.requireAuth(app.handleDeleteFile))
	mux.HandleFunc("POST /api/files/{id}/share", app.requireAuth(app.handleShareFile))

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
