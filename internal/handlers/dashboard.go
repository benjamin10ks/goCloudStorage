package handlers

import "net/http"

func TemplateRenderHandler(w http.ResponseWriter, r *http.Request) {
	// Render the dashboard template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><h1>Dashboard</h1></body></html>"))
}
