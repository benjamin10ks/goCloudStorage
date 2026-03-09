package main

import (
	"log"
	"net/http"
)

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := a.auth.Register(email, password)
	if err != nil {
		log.Printf("Registration failed for email %s: %v", email, err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}
	log.Printf("User registered: %s", user.Email)
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
}

// remove session cookie
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {}
