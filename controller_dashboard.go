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
		log.Printf("User not found in context")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	ctx := r.Context()

	files, err := getUserFiles(ctx, a.db, user.ID)
	if err != nil {
		log.Printf("Error fetching user files: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"User":  user,
		"Files": files,
	}

	err = a.tmpl["dashboard"].ExecuteTemplate(w, "app_layout", data)
	if err != nil {
		log.Printf("Error rendering dashboard: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
