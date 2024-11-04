package main

import (
	"log"
	"net/http"

	"github.com/ncfex/dcart/api-gateway/internal/routes"
)

// TODO - use config to set server port, service urls
func main() {
	router, err := routes.NewRouter("http://0.0.0.0:8081")
	if err != nil {
		log.Fatal(err)
	}

	handler := router.SetupRoutes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
