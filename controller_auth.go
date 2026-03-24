package main

import (
	"net/http"
)

// display landing page with login and registration options
func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// display login page with options for passkey and github login
func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *App) handlePasskeyBeginRegister(w http.ResponseWriter, r *http.Request)  {}
func (a *App) handlePasskeyFinishRegister(w http.ResponseWriter, r *http.Request) {}
func (a *App) handlePasskeyBeginLogin(w http.ResponseWriter, r *http.Request)     {}
func (a *App) handlePasskeyFinishLogin(w http.ResponseWriter, r *http.Request)    {}

// possibly extend in future to support more providers
func (a *App) handleGitHubLogin(w http.ResponseWriter, r *http.Request)    {}
func (a *App) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {}

// remove session cookie
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {}
