package main

import (
	"net/http"
	"viz_pbi-server/api"
)

func main() {
	srv := api.NewServer()
	http.ListenAndServe(":8080", srv)
}
