package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"viz_pbi-server/minio"
	"viz_pbi-server/mongo"
)

var fragmentsBucket string

func init() {
	fragmentsBucket = os.Getenv("MINIO_FRAGMENTS_BUCKET")
	if fragmentsBucket == "" {
		fragmentsBucket = "ifc-fragment-files"
	}
}

type Server struct {
	*mux.Router
	mongoReader *mongo.Reader
}

func NewServer() *Server {
	mongoUri := getEnv("MONGO_URI", "mongodb://localhost:27017")
	s := &Server{
		Router:      mux.NewRouter(),
		mongoReader: mongo.NewReader(mongoUri),
	}

	s.routes()

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	s.Router.Use(c.Handler)
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/file", s.getFile()).Methods("GET")
	s.HandleFunc("/data", s.getData()).Methods("GET")
}

func (s *Server) getFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
			return
		}

		// Check if the file name is valid
		passed, msg := checkFileName(name)
		if !passed {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Fetch the file from MinIO
		file, err := minio.GetFile(fragmentsBucket, name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch file from MinIO: %v", err), http.StatusInternalServerError)
			return
		}

		// Set the appropriate headers and return the file content
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(file)
	}
}

// gets all data in a mongodb collection, expect two parameters: db and collection
func (s *Server) getData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := r.URL.Query().Get("db")
		collection := r.URL.Query().Get("collection")
		if db == "" || collection == "" {
			http.Error(w, "Missing 'db' or 'collection' query parameter", http.StatusBadRequest)
			return
		}

		// Fetch the data from MongoDB
		data, err := s.mongoReader.ReadAllElements(db, collection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the appropriate headers and return the data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func checkFileName(name string) (bool, string) {
	// check if file ends with .gz
	if !strings.HasSuffix(name, ".gz") {
		return false, "file name does not end with .gz"
	}

	pattern := "project/filename_timestamp.gz"

	// check if file name follows the pattern: project/filename_timestamp.gz
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return false, pattern
	}

	// split the filename and timestamp
	filename := parts[1]
	filenameParts := strings.Split(filename, "_")
	if len(filenameParts) != 2 {
		return false, pattern
	}

	return true, "File name is valid"
}
