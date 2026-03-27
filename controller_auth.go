package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type PendingAuth struct {
	ExpiresAt time.Time
}

var (
	pendingAuths   = make(map[string]PendingAuth)
	pendingAuthsMu sync.Mutex
)

// display landing page with login and registration options
func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	a.tmpl.ExecuteTemplate(w, "register.tmpl", nil)
}

// display login page with options for passkey and github login
func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	a.tmpl.ExecuteTemplate(w, "login.tmpl", nil)
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

	err = saveWebAuthnSession(a.db, "registration", user.ID, sessionData)
	if err != nil {
		http.Error(w, "Failed to begin registration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func (a *App) handlePasskeyFinishRegister(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	user, err := getUserByID(a.db, userID)
	if err != nil {
		http.Error(w, "Failed to finish registration", http.StatusInternalServerError)
		return
	}

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

	sessionToken, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	err = saveSession(a.db, user.ID, sessionToken, "passkey")
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
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
	options, sessionData, err := a.auth.BeginLoginWithoutUser()
	if err != nil {
		http.Error(w, "Failed to begin login", http.StatusInternalServerError)
		return
	}
	err = saveWebAuthnSession(a.db, "login", 0, sessionData)
	if err != nil {
		http.Error(w, "Failed to begin login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func (a *App) handlePasskeyFinishLogin(w http.ResponseWriter, r *http.Request) {
	sessionData, err := getWebAuthnSession(a.db, "login", 0)
	if err != nil {
		http.Error(w, "Failed to finish login", http.StatusInternalServerError)
		return
	}

	user, credential, err := a.auth.FinishLoginWithoutUser(r, sessionData)
	if err != nil {
		http.Error(w, "Failed to finish login", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	err = updatePasskeySignCount(ctx, a.db, user.ID, credential.ID, credential.Authenticator.SignCount)
	if err != nil {
		http.Error(w, "Failed to update sign count", http.StatusInternalServerError)
		return
	}

	sessionToken, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}
	err = saveSession(a.db, user.ID, sessionToken, "passkey")
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (a *App) handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	stateToken, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate state token", http.StatusInternalServerError)
		return
	}

	pendingAuthsMu.Lock()
	pendingAuths[stateToken] = PendingAuth{
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	pendingAuthsMu.Unlock()

	authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=read:user&state=%s", GithubClientID, stateToken)

	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

func (a *App) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	pendingAuthsMu.Lock()
	pending, ok := pendingAuths[state]
	if ok {
		delete(pendingAuths, state)
	}
	pendingAuthsMu.Unlock()

	if !ok {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	if time.Now().After(pending.ExpiresAt) {
		http.Error(w, "Expired state parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	accessToken, err := exchangeCodeForToken(ctx, code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Error(w, "Error exchanging code for token", http.StatusInternalServerError)
		return
	}

	githubUserID, err := getGithubUserID(ctx, accessToken)
	if err != nil {
		http.Error(w, "Error fetching GitHub user info", http.StatusInternalServerError)
		return
	}

	user, err := getUserByGithubID(ctx, a.db, githubUserID)
	if err != nil {
		user, err = createUserFromGithub(ctx, a.db, githubUserID, accessToken)
		if err != nil {
			log.Printf("Error creating user with GitHub ID: %v", err)
			http.Error(w, "Error creating user with GitHub ID", http.StatusInternalServerError)
			return
		}
	}

	err = storeGithubToken(ctx, a.db, user.ID, githubUserID, accessToken)
	if err != nil {
		log.Printf("Error storing GitHub token: %v", err)
		http.Error(w, "Error storing GitHub token", http.StatusInternalServerError)
		return
	}

	sessionToken, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}
	err = saveSession(a.db, user.ID, sessionToken, "github")
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// remove session cookie
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ctx := r.Context()
	err = deleteSession(ctx, a.db, cookie.Value)
	if err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	log.Println("User logged out")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
