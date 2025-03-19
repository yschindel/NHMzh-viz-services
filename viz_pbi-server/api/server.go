package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"viz_pbi-server/minio"

	_ "github.com/marcboeker/go-duckdb"
)

var fragmentsBucket string
var lcaCostDataBucket string
var apiKey string

func init() {
	fragmentsBucket = getEnv("MINIO_FRAGMENTS_BUCKET", "ifc-fragment-files")
	lcaCostDataBucket = getEnv("MINIO_LCA_COST_DATA_BUCKET", "lca-cost-data")
	apiKey = getEnv("PBI_SRV_API_KEY", "")

	if apiKey == "" {
		log.Println("WARNING: No API key set. API endpoints are not secured!")
	}
}

type Server struct {
	http.Handler
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

func NewServer() *Server {
	// Create a new router
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/fragments", getFragmentsFile()).Methods("GET")
	router.HandleFunc("/fragments/list", getFragmentFileNames()).Methods("GET")
	router.HandleFunc("/data", getDataFile()).Methods("GET")
	router.HandleFunc("/data/list", getDataFileNames()).Methods("GET")

	// Create middleware chain
	var handler http.Handler = router

	// Add proxy headers middleware
	handler = proxyHeadersMiddleware(handler)

	// Apply API key middleware if API key is set
	if apiKey != "" {
		log.Println("API key authentication enabled")
		handler = apiKeyMiddleware(handler)
	}

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Apply CORS middleware
	handler = c.Handler(handler)

	return &Server{Handler: handler}
}

func getFragmentsFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
			return
		}

		log.Printf("getting fragments files with id: %s", id)

		files, err := minio.ListAllFiles(fragmentsBucket)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list all files from MinIO: %v", err), http.StatusInternalServerError)
			return
		}

		// find the file name that matches the id
		var fileNames []string
		for _, file := range files {
			if strings.Contains(file, id) {
				fileNames = append(fileNames, file)
			}
		}

		// get the file with the latest timestamp
		sort.Strings(fileNames)
		name := fileNames[len(fileNames)-1]

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

// get all file names in the ifc-fragment-files bucket
func getFragmentFileNames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := minio.ListAllFiles(fragmentsBucket)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list all files from MinIO: %v", err), http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			log.Printf("fragments file: %s", file)
		}

		// return list of files
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	}
}

// get a data file from the lca-cost-data bucket
func getDataFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
			return
		}

		project := r.URL.Query().Get("project")
		if project == "" {
			http.Error(w, "Missing 'project' query parameter", http.StatusBadRequest)
			return
		}

		fileName := project + "/" + name

		// Check if the file name is valid
		passed, msg := checkDataFileName(fileName)
		if !passed {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Fetch the file from MinIO
		file, err := minio.GetFile(lcaCostDataBucket, fileName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch file from MinIO: %v", err), http.StatusInternalServerError)
			return
		}

		// Set the appropriate headers and return the file content
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(file)
	}
}

// get all projects in the lca-cost-data bucket
func getDataFileNames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		files, err := minio.ListAllFiles(lcaCostDataBucket)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list all files from MinIO: %v", err), http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			log.Printf("data file: %s", file)
		}

		// return list of files
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func checkDataFileName(name string) (bool, string) {
	// check if file ends with .gz
	if !strings.HasSuffix(name, ".parquet") {
		return false, "file name does not end with .parquet"
	}

	pattern := "project/filename.parquet"

	// check if file name follows the pattern: project/filename.parquet
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return false, pattern
	}

	return true, "File name is valid"
}
