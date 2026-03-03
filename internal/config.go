// Package config contains server configuration.
package main

import (
	"net/http"
)

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
