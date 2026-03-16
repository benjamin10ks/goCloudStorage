package main

import (
	"log"
	"net/http"
)

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	user, err := a.auth.Register(token, a.db)
	if err != nil {
		log.Printf("Registration failed for email %s: %v", email, err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}
	log.Printf("User registered: %s", user.Email)
	w.WriteHeader(http.StatusCreated)
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	err := a.auth.Login(token)
	if err != nil {
		log.Printf("Login failed: %v", err)
	}
}

// remove session cookie
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {}
