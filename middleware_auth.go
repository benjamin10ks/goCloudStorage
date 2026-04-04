package main

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const contextKeyUser contextKey = "user"

func (a *App) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Printf("No session cookie found: %v", err)
			http.Error(w, "Unauthorized: no session cookie", http.StatusUnauthorized)
			return
		}

		user, err := getUserBySessionToken(r.Context(), a.db, cookie.Value)
		if err != nil {
			log.Printf("Failed to get user from session: %v", err)
			http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next(w, r.WithContext(ctx))
	}
}

func userFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(contextKeyUser).(*User)
	return user, ok
}
