package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	kafkaBroker     = "localhost:9092" // has to be hardcoded for the test because the .env uses 'docker-compose' service name
	topic           string
	ifcBucket       string
	fragmentsBucket string
	ifcApiPort      int
	pbiServerPort   int
	pbiServerUrl    string
	messages        []Message
)

func init() {
	// Load environment variables from .env file
	envPath := filepath.Join("..", "..", ".env")
	log.Printf("Loading environment variables from %s", envPath)
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Read environment variables

	topic = getEnv("KAFKA_IFC_TOPIC", "ifc-files")
	ifcBucket = getEnv("MINIO_IFC_BUCKET", "ifc-files")
	fragmentsBucket = getEnv("MINIO_FRAGMENTS_BUCKET", "ifc-fragment-files")
	ifcApiPort = getEnvAsInt("IFC_API_PORT", 4242)
	pbiServerPort = getEnvAsInt("PBI_SERVER_PORT", 3000)
	pbiServerUrl = "http://localhost:" + strconv.Itoa(pbiServerPort)

	// log all environment variables
	log.Printf("KAFKA_BROKER: %s", kafkaBroker)
	log.Printf("KAFKA_IFC_TOPIC: %s", topic)
	log.Printf("MINIO_IFC_BUCKET: %s", ifcBucket)
	log.Printf("IFC_API_PORT: %d", ifcApiPort)
	log.Printf("MINIO_FRAGMENTS_BUCKET: %s", fragmentsBucket)
	log.Printf("PBI_SERVER_PORT: %d", pbiServerPort)

	// Initialize messages
	messages = []Message{
		newMessage("project1", "file1.ifc"),
		newMessage("project2", "file2.ifc"),
		newMessage("project1", "file3.ifc"),
	}

}

type Message struct {
	Project   string `json:"project"`
	Filename  string `json:"filename"`
	Location  string `json:"location"`
	Timestamp string `json:"timestamp"`
}

func newMessage(project, filename string) Message {
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	log.Println(timestamp)
	return Message{
		Project:   project,
		Filename:  filename,
		Location:  newFilename(project, filename, timestamp),
		Timestamp: timestamp,
	}
}

func newFilename(project, filename, timestamp string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	fileTimestamp := timestamp
	return fmt.Sprintf("%s/%s_%s%s", project, name, fileTimestamp, ext)
}

func addIfcFilesToMinio() error {
	// Read the test file content
	content, err := os.ReadFile(filepath.Join("..", "assets", "test.ifc"))
	if err != nil {
		return fmt.Errorf("failed to read test.ifc file: %w", err)
	}

	// Create a new HTTP client
	client := &http.Client{}

	for _, msg := range messages {
		// Create a new multipart form buffer
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add the file
		part, err := writer.CreateFormFile("file", msg.Filename)
		if err != nil {
			return fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := part.Write(content); err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}

		// Add the project field
		if err := writer.WriteField("project", msg.Project); err != nil {
			return fmt.Errorf("failed to write project field: %w", err)
		}

		// Add the timestamp field
		if err := writer.WriteField("timestamp", msg.Timestamp); err != nil {
			return fmt.Errorf("failed to write timestamp field: %w", err)
		}

		// Close the writer
		writer.Close()

		// Create the request
		url := fmt.Sprintf("http://localhost:%d/api/upload", ifcApiPort)
		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set the content type
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		// Check the response
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
		}

		// Parse the response to get the location
		var response struct {
			Location string `json:"location"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		// Update the message location with the one returned from the API
		log.Printf("Uploaded file to API: %s\n", response.Location)
	}

	return nil
}

func getFileFromPbiServer(url string, location string) ([]byte, error) {
	// Make HTTP request to /file endpoint
	resp, err := http.Get(url + "/fragments?name=" + location)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to PBI server: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from PBI server: got %v, want %v.\nError: %s", resp.StatusCode, http.StatusOK, err)
	}

	// Read the response body
	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return fileContent, nil
}

func main() {
	log.Println("Adding IFC files to MinIO...")
	if err := addIfcFilesToMinio(); err != nil {
		log.Fatalf("Error adding IFC files to MinIO: %v", err)
	}

	log.Println("Waiting for IFC consumer to convert ifc to gz and write back to MinIO...")
	time.Sleep(30 * time.Second)

	log.Println("Getting fragments files...")
	for _, msg := range messages {
		name := strings.Replace(msg.Location, ".ifc", ".gz", 1)
		_, err := getFileFromPbiServer(pbiServerUrl, name)
		if err != nil {
			log.Fatalf("Error getting fragments files: %v", err)
		}

		log.Printf("Received fragments file: %s\n", msg.Location)
	}

	log.Println("IFC test completed successfully")
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not set, using default value", key)
		return defaultValue
	}
	return value
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
