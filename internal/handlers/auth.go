package handlers

import (
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

func (h *Handler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := h.UserService.CreateUser(username, password)
	if err != nil {
		log.Printf("failed to create user, %d", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func HandleLogout() {
	// placeholder implementation
}
