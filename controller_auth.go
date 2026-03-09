package main

import "net/http"

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := register(email, password)
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
}

// remove session cookie
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {}
