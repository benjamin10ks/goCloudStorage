package main

import "net/http"

func (a *App) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (a *App) handleHomePage(w http.ResponseWriter, r *http.Request) {}
