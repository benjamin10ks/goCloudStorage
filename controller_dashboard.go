package main

import (
	"log"
	"net/http"
)

func (a *App) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	if a.auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	}
}

func (a *App) handleHomePage(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
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

	err := a.tmpl["dashboard"].ExecuteTemplate(w, "app_layout", data)
	if err != nil {
		log.Printf("Error rendering dashboard: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
