// Package handlers holds all main routes handlers
package handlers

import "net/http"

func PlaceHolderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is a placeholder handler."))
}
