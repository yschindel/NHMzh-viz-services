package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"viz_pbi-server/env"
	"viz_pbi-server/logger"
	"viz_pbi-server/storage"
	"viz_pbi-server/storage/azure"
)

var apiKey string
var log = logger.WithFields(logger.Fields{"component": "api"})

func init() {
	apiKey = env.Get("PBI_SRV_API_KEY")
}

type Server struct {
	http.Handler
	storage storage.StorageProvider
	logger  *logger.Logger
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
		log.Debug("Proxy headers detected: proto=%s, host=%s",
			r.Header.Get("X-Forwarded-Proto"),
			r.Header.Get("X-Forwarded-Host"))
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
			log.Warn("Unauthorized access attempt: remote=%s, path=%s", r.RemoteAddr, r.URL.Path)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid or missing API key"})
			return
		} else {
			log.Debug("API key authentication successful: remote=%s, path=%s", r.RemoteAddr, r.URL.Path)
		}

		// Call the next handler if API key is valid
		next.ServeHTTP(w, r)
	})
}

// Request logging middleware
func (s *Server) requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("Incoming request")
		next.ServeHTTP(w, r)
	})
}

func NewServer() (*Server, error) {
	log.Info("Initializing server")

	// Initialize Azure storage
	config := azure.NewConfig()
	storageProvider, err := azure.NewBlobStorage(config)
	if err != nil {
		log.Error("Failed to initialize Azure storage: %v", err)
		return nil, fmt.Errorf("failed to initialize Azure storage: %v", err)
	}

	server := &Server{
		storage: storageProvider,
		Handler: nil,
		logger:  logger.New(),
	}

	// Create router
	router := mux.NewRouter()

	// Register routes with query parameters for container
	router.HandleFunc("/files", server.handleGetFile()).Methods("GET")
	router.HandleFunc("/files", server.handleUploadFile()).Methods("POST")

	// Add middleware
	var handler http.Handler = router
	handler = server.requestLoggingMiddleware(handler) // Use server's logger
	handler = proxyHeadersMiddleware(handler)
	if apiKey != "" {
		log.Info("API key authentication enabled")
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

	log.Info("Server initialized successfully")
	return server, nil
}

func (s *Server) handleGetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get container and file from query parameters
		container := r.URL.Query().Get("container")
		if container == "" {
			container = s.storage.Container() // Use default container if not specified
		}
		file := r.URL.Query().Get("file")
		if file == "" {
			http.Error(w, "file parameter is required", http.StatusBadRequest)
			return
		}

		s.logger.Debug("Handling file download request: container=%s, file=%s",
			container, file)

		fileData, err := s.storage.GetFile(r.Context(), container, file)
		if err != nil {
			s.logger.Error("Failed to fetch file: container=%s, file=%s, error=%v",
				container, file, err)
			http.Error(w, fmt.Sprintf("Failed to fetch file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileData)

		s.logger.Info("File downloaded successfully: container=%s, file=%s",
			container, file)
	}
}

func (s *Server) handleUploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get container and file from query parameters
		container := r.URL.Query().Get("container")
		if container == "" {
			container = s.storage.Container() // Use default container if not specified
		}
		file := r.URL.Query().Get("file")
		if file == "" {
			http.Error(w, "file parameter is required", http.StatusBadRequest)
			return
		}

		s.logger.Debug("Handling file upload request: container=%s, file=%s",
			container, file)

		// Limit the file size to 1GB, just for safety
		r.Body = http.MaxBytesReader(w, r.Body, 1<<30)

		err := s.storage.UploadFile(r.Context(), container, file, r.Body)
		if err != nil {
			s.logger.Error("Failed to upload file: container=%s, file=%s, error=%v",
				container, file, err)
			http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
			return
		}

		s.logger.Info("File uploaded successfully: container=%s, file=%s",
			container, file)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":   "File uploaded successfully",
			"file":      file,
			"container": container,
		})
	}
}
