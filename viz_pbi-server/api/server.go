package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"viz_pbi-server/minio"
)

var fragmentsBucket string
var lcaCostDataBucket string

func init() {
	fragmentsBucket = getEnv("MINIO_FRAGMENTS_BUCKET", "ifc-fragment-files")
	lcaCostDataBucket = getEnv("MINIO_LCA_COST_DATA_BUCKET", "lca-cost-data")
}

type Server struct {
	*mux.Router
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
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
	s.HandleFunc("/file", s.getFragmentsFile()).Methods("GET")
	s.HandleFunc("/data/models", s.getDataFileList()).Methods("GET")
	s.HandleFunc("/data/projects", s.getProjects()).Methods("GET")
	s.HandleFunc("/data", s.getDataFile()).Methods("GET")
}

func (s *Server) getFragmentsFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
			return
		}

		// Check if the file name is valid
		passed, msg := checkFragmentsFileName(name)
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

// get a data file from the lca-cost-data bucket
func (s *Server) getDataFile() http.HandlerFunc {
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

// get all file names in the specified directory in the lca-cost-data bucket
func (s *Server) getDataFileList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		if project == "" {
			http.Error(w, "Missing 'project' query parameter", http.StatusBadRequest)
			return
		}

		files, err := minio.ListFiles(lcaCostDataBucket, project)
		log.Printf("data files: %v", files)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list files from MinIO: %v", err), http.StatusInternalServerError)
			return
		}

		// Filter only .parquet files and remove the project prefix
		var modelList []string
		for _, file := range files {
			if strings.HasSuffix(file, ".parquet") {
				// Remove the project prefix (e.g., "project1/")
				modelName := strings.TrimPrefix(file, project+"/")
				modelList = append(modelList, modelName)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(modelList); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// get all projects in the lca-cost-data bucket
func (s *Server) getProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projects, err := minio.ListProjects(lcaCostDataBucket)
		log.Printf("projects: %v", projects)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list projects from MinIO: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func checkFragmentsFileName(name string) (bool, string) {
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
