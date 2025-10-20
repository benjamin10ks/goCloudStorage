package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/benjamin10ks/goCloudStorage/internal/services"
)

type Handler struct {
	UserService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}

func HandleLogin() {
	// placeholder implementation
}

// change to respond with HTML
func (h *Handler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := h.UserService.CreateUser(username, password)
	if err != nil {
		log.Printf("failed to create user, %d", err)
	}
	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to marshal user data, %d", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func HandleLogout() {
	// placeholder implementation
}
