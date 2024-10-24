package main

import (
	"net/http"
	"os"
	"viz_pbi-server/api"
)

func main() {
	port := os.Getenv("PBI_SERVER_PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}
	srv := api.NewServer()
	http.ListenAndServe(":"+port, srv)
}
