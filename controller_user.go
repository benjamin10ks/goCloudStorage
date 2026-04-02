package main

import (
	"log"
	"net/http"
)

func (a *App) handleUserProfile(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		log.Printf("User not found in context")
		http.Error(w, "Unauthorized: user not found in context", http.StatusUnauthorized)
		return
	}

	NewDiplayName := r.FormValue("display_name")

	ctx := r.Context()

	err := updateUserDisplayName(ctx, a.db, user.ID, NewDiplayName)
	if err != nil {
		log.Printf("Error updating display name: %v", err)
		http.Error(w, "Failed to update display name", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (a *App) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleListUserFiles(w http.ResponseWriter, r *http.Request) {
}
