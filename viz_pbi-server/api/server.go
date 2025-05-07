package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"viz_pbi-server/env"
	"viz_pbi-server/logger"
	"viz_pbi-server/models"
	"viz_pbi-server/storage"
	"viz_pbi-server/storage/azure"
)

var apiKey string
var log = logger.WithFields(logger.Fields{"component": "api/server.go"})

func init() {
	apiKey = env.Get("API_KEY")
}

type Server struct {
	http.Handler
	storage   storage.StorageProvider
	logger    *logger.Logger
	db        *sql.DB
	sqlWriter *azure.SqlWriter
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
		reqLogger := log.WithFields(logger.Fields{
			"proto": r.Header.Get("X-Forwarded-Proto"),
			"host":  r.Header.Get("X-Forwarded-Host"),
		})
		reqLogger.Debug("Proxy headers detected")
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
			reqLogger := log.WithFields(logger.Fields{
				"remote_addr": r.RemoteAddr,
				"path":        r.URL.Path,
			})
			reqLogger.Warn("Unauthorized access attempt")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid or missing API key"})
			return
		} else {
			reqLogger := log.WithFields(logger.Fields{
				"remote_addr": r.RemoteAddr,
				"path":        r.URL.Path,
			})
			reqLogger.Debug("API key authentication successful")
		}

		// Call the next handler if API key is valid
		next.ServeHTTP(w, r)
	})
}

// Request logging middleware
func (s *Server) requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := s.logger.WithFields(logger.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		})
		reqLogger.Debug("Incoming request")
		next.ServeHTTP(w, r)
	})
}

func NewServer() (*Server, error) {
	log.Info("Initializing server")

	// Initialize Azure storage
	config := azure.NewConfig()
	storageProvider, err := azure.NewBlobStorage(config)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Failed to initialize Azure storage")
		return nil, fmt.Errorf("failed to initialize Azure storage: %v", err)
	}

	dbConfig := azure.DBConfig{
		Server:   env.Get("AZURE_DB_SERVER"),
		Port:     StringToInt(env.Get("AZURE_DB_PORT")),
		User:     env.Get("AZURE_DB_USER"),
		Password: env.Get("AZURE_DB_PASSWORD"),
		Database: env.Get("AZURE_DB_DATABASE"),
	}

	db, err := azure.ConnectDB(dbConfig)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Failed to connect to Azure SQL database")
		return nil, fmt.Errorf("failed to connect to Azure SQL database: %v", err)
	}

	// Initialize database schema
	err = azure.InitializeDatabase(db)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Failed to initialize database schema")
		return nil, fmt.Errorf("failed to initialize database schema: %v", err)
	}

	server := &Server{
		storage:   storageProvider,
		Handler:   nil,
		logger:    logger.New(),
		db:        db,
		sqlWriter: azure.NewSqlWriter(db),
	}

	// Create router
	router := mux.NewRouter()

	blobEndpoint := env.Get("STORAGE_ENDPOINT_FILE")
	if !strings.HasPrefix(blobEndpoint, "/") {
		blobEndpoint = "/" + blobEndpoint
	}

	materialsEndpoint := env.Get("STORAGE_ENDPOINT_DATA_MATERIALS")
	if !strings.HasPrefix(materialsEndpoint, "/") {
		materialsEndpoint = "/" + materialsEndpoint
	}

	elementsEndpoint := env.Get("STORAGE_ENDPOINT_DATA_ELEMENTS")
	if !strings.HasPrefix(elementsEndpoint, "/") {
		elementsEndpoint = "/" + elementsEndpoint
	}

	// Register routes with query parameters for container
	router.HandleFunc(blobEndpoint, server.handleGetBlob()).Methods("GET")
	router.HandleFunc(blobEndpoint, server.handleUploadBlob()).Methods("POST")
	router.HandleFunc(materialsEndpoint, server.handlePostMaterialData()).Methods("POST")
	router.HandleFunc(elementsEndpoint, server.handlePostElementsData()).Methods("POST")

	// Add middleware
	var handler http.Handler = router
	if apiKey != "" {
		handler = apiKeyMiddleware(handler)
		log.Info("API key authentication enabled")
	}
	handler = proxyHeadersMiddleware(handler)
	handler = server.requestLoggingMiddleware(handler) // Use server's logger

	// Enable CORS
	// This is required by PowerBI
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*", "null"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// AllowCredentials: true,
		ExposedHeaders: []string{
			"X-Metadata-Filename",
			"X-Metadata-Projectname",
			"X-Metadata-Timestamp",
			"Content-Type",
			"Content-Length",
		},
		MaxAge: 86400,
		Debug:  true,
	})

	handler = c.Handler(handler)
	server.Handler = handler

	log.Info("Server initialized successfully")
	return server, nil
}

func (s *Server) handleGetBlob() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get container and file from query parameters
		reqLogger := s.logger.WithFields(logger.Fields{
			"endpoint": "getBlob",
		})

		reqLogger.Debug("Handling file download request")

		container := r.URL.Query().Get("container")
		if container == "" {
			container = s.storage.Container() // Use default container if not specified
		}
		blobID := r.URL.Query().Get("id")
		if blobID == "" {
			http.Error(w, "id parameter is required", http.StatusBadRequest)
			return
		}

		reqLogger.WithFields(logger.Fields{
			"container": container,
			"id":        blobID,
		}).Debug("Url Parameters")

		fileData, err := s.storage.GetBlob(r.Context(), container, blobID)
		if err != nil {
			reqLogger.WithFields(logger.Fields{"error": err}).Error("Failed to fetch file")
			http.Error(w, fmt.Sprintf("Failed to fetch file: %v", err), http.StatusInternalServerError)
			return
		}

		reqLogger.Debug("File fetched successfully")

		// get the metadata from the blob
		metadata, err := s.storage.GetBlobMetadata(r.Context(), container, blobID)
		if err != nil {
			reqLogger.WithFields(logger.Fields{"error": err}).Error("Failed to fetch file metadata")
			http.Error(w, fmt.Sprintf("Failed to fetch file metadata: %v", err), http.StatusInternalServerError)
			return
		}

		// log the metadata
		reqLogger.WithFields(logger.Fields{"metadata": metadata}).Debug("Fetched File Metadata:")

		w.Header().Set("Content-Type", "application/octet-stream")

		// set the metadata as headers
		for key, value := range metadata {
			if value != nil { // Check if pointer is not nil
				headerKey := fmt.Sprintf("X-Metadata-%s", key)
				w.Header().Set(headerKey, *value) // Dereference the pointer
				reqLogger.WithFields(logger.Fields{
					"header_key":   headerKey,
					"header_value": *value,
				}).Debug("Setting metadata header")
			}
		}

		w.Write(fileData)

		reqLogger.Info("File downloaded successfully")
	}
}

func (s *Server) handleUploadBlob() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit the file size to 1GB
		if err := r.ParseMultipartForm(1 << 30); err != nil {
			s.logger.WithFields(logger.Fields{"error": err}).Error("Failed to parse form")
			http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
			return
		}
		defer r.MultipartForm.RemoveAll()

		// Get the file
		file, _, err := r.FormFile("file")
		if err != nil {
			s.logger.WithFields(logger.Fields{"error": err}).Error("File is required")
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// use the fileID from the form and internally call the varialbe blobID
		// to match the azure blob storage naming convention
		blobID := r.FormValue("fileID")
		if blobID == "" {
			s.logger.Error("fileID is required")
			http.Error(w, "fileID is required", http.StatusBadRequest)
			return
		}

		// Get metadata from form fields
		// this will be used to set the blob metadata
		fileName := r.FormValue("fileName")
		if fileName == "" {
			s.logger.Error("fileName is required")
			http.Error(w, "fileName (original filename) is required", http.StatusBadRequest)
			return
		}

		projectName := r.FormValue("projectName")
		if projectName == "" {
			s.logger.Error("projectName is required")
			http.Error(w, "projectName is required", http.StatusBadRequest)
			return
		}

		timestamp := r.FormValue("timestamp")
		if timestamp == "" {
			s.logger.Error("timestamp is required")
			http.Error(w, "timestamp is required", http.StatusBadRequest)
			return
		}

		// Get container from form or use default
		container := r.FormValue("container")
		if container == "" {
			container = s.storage.Container()
		}

		blobData := models.BlobData{
			Container: container,
			Project:   projectName,
			Filename:  fileName,
			Timestamp: timestamp,
			BlobID:    blobID,
			Blob:      file,
		}

		// Create a logger with upload context
		uploadLogger := s.logger.WithFields(logger.Fields{
			"container":   blobData.Container,
			"fileName":    blobData.Filename,
			"projectName": blobData.Project,
			"timestamp":   blobData.Timestamp,
		})

		uploadLogger.Debug("Starting file upload")

		blobData, err = s.storage.UploadBlob(r.Context(), blobData)
		if err != nil {
			uploadLogger.WithFields(logger.Fields{"error": err}).Error("Failed to upload file")
			http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
			return
		}

		// Update the data_updates table with the new blobID
		err = s.sqlWriter.WriteBlobData(blobData)
		if err != nil {
			uploadLogger.WithFields(logger.Fields{"error": err}).Error("Failed to write blob upload message")
			http.Error(w, fmt.Sprintf("Failed to write blob upload message: %v", err), http.StatusInternalServerError)
			return
		}

		uploadLogger.Info("File uploaded successfully")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":   "File uploaded successfully",
			"blobID":    blobID,
			"container": container,
		})
	}
}

func (s *Server) handlePostMaterialData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLogger := s.logger.WithFields(logger.Fields{
			"endpoint": "materials",
		})

		var materialData []models.EavMaterialDataItem
		if err := json.NewDecoder(r.Body).Decode(&materialData); err != nil {
			reqLogger.WithFields(logger.Fields{"error": err}).Error("Failed to decode request body")
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		if err := s.sqlWriter.WriteMaterials(materialData); err != nil {
			reqLogger.WithFields(logger.Fields{"error": err}).Error("Failed to write LCA data")
			http.Error(w, fmt.Sprintf("Failed to write LCA data: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Materials data written successfully",
		})
	}
}

func (s *Server) handlePostElementsData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLogger := s.logger.WithFields(logger.Fields{
			"endpoint": "elements",
		})

		var elementsData []models.EavElementDataItem
		if err := json.NewDecoder(r.Body).Decode(&elementsData); err != nil {
			reqLogger.WithFields(logger.Fields{"error": err}).Error("Failed to decode request body")
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		if err := s.sqlWriter.WriteElements(elementsData); err != nil {
			reqLogger.WithFields(logger.Fields{"error": err}).Error("Failed to write elements data")
			http.Error(w, fmt.Sprintf("Failed to write elements data: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Elements data written successfully",
		})
	}
}
