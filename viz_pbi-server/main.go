package main

import (
	"net/http"
	"viz_pbi-server/api"
)

func main() {
	srv := api.NewServer()
	http.ListenAndServe(":3000", srv)
}
