package main

import (
	"net/http"
)

const (
	BASE_PATH = "./storage/users/"
)

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
