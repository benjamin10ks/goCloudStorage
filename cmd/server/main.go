package main

import (
	"fmt"

	"github.com/benjamin10ks/goCloudStorage/internal/config"
)

func main() {
	fmt.Printf("Starting server on port %s \n", config.Server.Addr)
	config.Server.ListenAndServe()
}
