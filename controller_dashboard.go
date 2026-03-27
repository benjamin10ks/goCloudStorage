package main

import "net/http"

func (a *App) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	if a.auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (a *App) handleHomePage(w http.ResponseWriter, r *http.Request) {
	// a.tmpl["dashboard.tmpl"].ExecuteTemplate(w, "home.tmpl", nil)
}
