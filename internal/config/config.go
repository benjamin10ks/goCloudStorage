// Package config contains server configuration.
package config

import (
	"net/http"
)

func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
