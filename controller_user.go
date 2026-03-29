package main

import (
	"log"
	"net/http"
)

func (a *App) handleUserProfile(w http.ResponseWriter, r *http.Request) {
}

func (a *App) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	NewDiplayName := r.FormValue("display_name")
	sessionToken, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("Error retrieving session cookie: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	user, err := getUserBySessionToken(ctx, a.db, sessionToken.Value)
	if err != nil {
		log.Printf("Error getting user from session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = updateUserDisplayName(ctx, a.db, user.ID, NewDiplayName)
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
