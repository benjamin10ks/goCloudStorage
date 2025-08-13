// Package config contains server configuration.
package config

import (
	"net/http"
)

var mux = http.NewServeMux()

var Server = &http.Server{
	Addr:    ":8080",
	Handler: mux,
}
