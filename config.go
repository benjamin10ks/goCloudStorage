package main

import (
	"net/http"
	"os"
)

const BASE_PATH = "./storage/users/"

// IMPORTANT CONSTANTS
var (
	GithubClientID = os.Getenv("GITHUB_CLIENT_ID")
	GithubSecret   = os.Getenv("GITHUB_CLIENT_SECRET")
)

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
