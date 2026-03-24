package main

import (
	"net/http"
	"os"
)

// IMPORTANT CONSTANTS
var (
	BASE_PATH      = "./storage/users/"
	GithubClientID = os.Getenv("GITHUB_CLIENT_ID")
	GithubSecret   = os.Getenv("GITHUB_CLIENT_SECRET")
)

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
