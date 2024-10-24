package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"viz_pbi-server/minio"
)

type Server struct {
	*mux.Router
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/files/getFirst", s.getFile()).Methods("GET")
}

func (s *Server) getFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
			return
		}

		// Fetch the file from MinIO
		bucketName := "ifc-fragment-files" // Replace with your actual bucket name
		fileContent, err := minio.GetFile(bucketName, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the appropriate headers and return the file content
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileContent)
	}
}
