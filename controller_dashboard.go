package main

import (
	"log"
	"net/http"
)

func (a *App) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	if a.auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (a *App) handleHomePage(w http.ResponseWriter, r *http.Request) {
	if !a.auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	
	// Get user from session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	
	ctx := r.Context()
	user, err := getUserBySessionToken(ctx, a.db, cookie.Value)
	if err != nil {
		log.Printf("Error getting user from session: %v", err)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	
	// Get user's files
	// TODO: Implement getUserFiles function
	// files, err := getUserFiles(a.db, user.ID)
	
	data := map[string]any{
		"User": user,
		// "Files": files,
	}
	
	err = a.tmpl["dashboard"].ExecuteTemplate(w, "app_layout", data)
	if err != nil {
		log.Printf("Error rendering dashboard: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
