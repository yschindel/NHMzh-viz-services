package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"viz_pbi-server/storage"
	"viz_pbi-server/storage/azure"
)

var apiKey string

func init() {
	apiKey = getEnv("PBI_SRV_API_KEY", "")

	if apiKey == "" {
		log.Println("WARNING: No API key set. API endpoints are not secured!")
	}
}

type Server struct {
	http.Handler
	storage storage.StorageProvider
}

// Proxy-aware middleware to handle X-Forwarded headers
func proxyHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If X-Forwarded-Proto is set and is https, update the request scheme
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			r.URL.Scheme = "https"
		}

		// If X-Forwarded-Host is set, update the request host
		if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
			r.Host = forwardedHost
		}

		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}

// API key middleware to secure all endpoints
func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get API key from request header
		key := r.Header.Get("X-API-Key")

		// Check if API key is valid
		if key == "" || key != apiKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid or missing API key"})
			return
		}

		// Call the next handler if API key is valid
		next.ServeHTTP(w, r)
	})
}

func NewServer() (*Server, error) {
	// Initialize Azure storage
	config := azure.NewConfig()
	storageProvider, err := azure.NewBlobStorage(config)
	if err != nil {
		return nil, err
	}

	server := &Server{
		storage: storageProvider,
		Handler: nil,
	}

	// Create router
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/files/{container}/{filename}", server.handleGetFile()).Methods("GET")
	router.HandleFunc("/files/{container}/{filename}", server.handleUploadFile()).Methods("POST")

	// Add middleware
	var handler http.Handler = router
	handler = proxyHeadersMiddleware(handler)
	if apiKey != "" {
		log.Println("API key authentication enabled")
		handler = apiKeyMiddleware(handler)
	}

	// Enable CORS
	// This is required by PowerBI
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler = c.Handler(handler)
	server.Handler = handler

	return server, nil
}

func (s *Server) handleGetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		container := vars["container"]
		filename := vars["filename"]

		file, err := s.storage.GetFile(r.Context(), container, filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(file)
	}
}

func (s *Server) handleUploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		container := vars["container"]
		filename := vars["filename"]

		// Limit the file size to 1GB, just for safety
		r.Body = http.MaxBytesReader(w, r.Body, 1<<30)

		err := s.storage.UploadFile(r.Context(), container, filename, r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "File uploaded successfully",
			"file":    filename,
		})
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
