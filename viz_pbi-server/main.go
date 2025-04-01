package main

import (
	"log"
	"net/http"
	"os"
	"viz_pbi-server/api"
)

func main() {
	port := os.Getenv("PBI_SERVER_PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}

	log.Printf("Starting server on port %s", port)
	srv, err := api.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Server is designed to run behind a reverse proxy for HTTPS termination")
	log.Fatal(http.ListenAndServe(":"+port, srv.Handler))
}
