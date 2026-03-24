package main

import (
	"encoding/json"
	"net/http"
)

// display landing page with login and registration options
func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Registration Page - choose Passkey or GitHub"))
	w.WriteHeader(http.StatusOK)
}

// display login page with options for passkey and github login
func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login Page - choose Passkey or GitHub"))
	w.WriteHeader(http.StatusOK)
}

// These use javascript client side to handle passkey flow
func (a *App) handleRegisterUsername(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	userID, err := createUser(a.db, username)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		a.tmpl.ExecuteTemplate(w, "error.tmpl", "Username already exists or invalid")
		return
	}
	a.tmpl.ExecuteTemplate(w, "passkey_begin.tmpl", map[string]any{
		"UserID": userID,
	})
}

func (a *App) handlePasskeyBeginRegister(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	user, err := getUserByID(a.db, userID)
	if err != nil {
		http.Error(w, "Failed to begin registration", http.StatusInternalServerError)
		return
	}
	options, sessionData, err := a.auth.BeginRegistration(user)
	if err != nil {
		http.Error(w, "Failed to begin registration", http.StatusInternalServerError)
		return
	}

	saveWebAuthnSession(a.db, "registration", user.ID, sessionData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func (a *App) handlePasskeyFinishRegister(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	user, err := getUserByID(a.db, userID)

	sessionData, err := getWebAuthnSession(a.db, "registration", user.ID)
	if err != nil {
		http.Error(w, "Failed to finish registration", http.StatusInternalServerError)
		return
	}

	if err := a.auth.FinishRegistration(user, sessionData, r); err != nil {
		http.Error(w, "Failed to finish registration", http.StatusInternalServerError)
		a.tmpl.ExecuteTemplate(w, "error.tmpl", "Passkey registration failed. Please try again.")
		return
	}

	sessionToken := generateToken()
	saveSession(a.db, user.ID, sessionToken, "passkey")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (a *App) handlePasskeyBeginLogin(w http.ResponseWriter, r *http.Request) {
	options, sessionData, err := a.auth.BeginLoginWitoutUser()
	if err != nil {
		http.Error(w, "Failed to begin login", http.StatusInternalServerError)
		return
	}
	saveWebAuthnSession(a.db, "login", 0, sessionData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}
func (a *App) handlePasskeyFinishLogin(w http.ResponseWriter, r *http.Request) {}

// possibly extend in future to support more providers
func (a *App) handleGitHubLogin(w http.ResponseWriter, r *http.Request)    {}
func (a *App) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {}

// remove session cookie
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {}
