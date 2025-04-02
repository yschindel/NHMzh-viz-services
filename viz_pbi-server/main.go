package main

import (
	"net/http"
	"os"
	"viz_pbi-server/api"
	"viz_pbi-server/logger"
)

func main() {
	log := logger.New()

	port := os.Getenv("PBI_SERVER_PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}

	log.Info("Starting server on port %s", port)

	srv, err := api.NewServer()
	if err != nil {
		log.Error("Failed to create server: %v", err)
		os.Exit(1)
	}

	log.Info("Server is designed to run behind a reverse proxy for HTTPS termination")

	if err := http.ListenAndServe(":"+port, srv.Handler); err != nil {
		log.Error("Server failed: %v", err)
		os.Exit(1)
	}

}
