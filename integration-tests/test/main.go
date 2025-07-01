// Package main implements integration tests for IFC file processing.
// It is designed to perform a full end-to-end test of the IFC file processing pipeline.
// It adds IFC files to MinIO directly and sends them to the Kafka topic.
// Then it waits for the IFC consumer to convert the IFC files to fragments and save them to MinIO.
// Finally it gets the fragments files from the PBI server and verifies the processing pipeline.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaBroker   = "localhost:9092" // has to be hardcoded for the test because the .env uses 'docker-compose' service name
	topic         string
	ifcBucket     string
	pbiServerPort int
	pbiServerUrl  string
	messages      []TestFileData
	minioClient   *minio.Client
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
	pbiServerPort = getEnvAsInt("PBI_SERVER_PORT", 3000)
	pbiServerUrl = "http://localhost:" + strconv.Itoa(pbiServerPort)

	// log all environment variables
	log.Printf("KAFKA_BROKER: %s", kafkaBroker)
	log.Printf("KAFKA_IFC_TOPIC: %s", topic)
	log.Printf("MINIO_IFC_BUCKET: %s", ifcBucket)
	log.Printf("PBI_SERVER_PORT: %d", pbiServerPort)

	// Initialize messages
	messages = []TestFileData{
		newTestFileData("project1", "file1.ifc"),
		// newTestFileData("project2", "file2.ifc"),
		// newTestFileData("project1", "file3.ifc"),
	}

	// Initialize MinIO client
	minioEndpoint := "localhost:9000"
	accessKeyID := getEnv("MINIO_ACCESS_KEY", "ROOTUSER")
	secretAccessKey := getEnv("MINIO_SECRET_KEY", "CHANGEME123")
	useSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	// Split endpoint and port if needed
	endpoint := minioEndpoint
	if !strings.Contains(endpoint, ":") {
		port := "9000"
		endpoint = fmt.Sprintf("%s:%s", endpoint, port)
	}

	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Error initializing MinIO client: %v", err)
	}
}

type TestFileData struct {
	Project          string
	FilenameOriginal string
	FileNameUUID     string
	Timestamp        string
}

func newTestFileData(project, filename string) TestFileData {
	// timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	timestamp := "2025-05-07T16:46:06.999Z"
	log.Println(timestamp)
	fileUUID := uuid.New().String()
	return TestFileData{
		Project:          project,
		FilenameOriginal: filename,
		FileNameUUID:     fileUUID + ".ifc",
		Timestamp:        timestamp,
	}
}

// Minio client configuration is already defined, but let's ensure the MinIO bucket exists
func ensureMinIOBucketExists() error {
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, ifcBucket)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, ifcBucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket: %s", ifcBucket)
	}

	return nil
}

// Add Kafka helper function to create a writer
func newKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func addIfcFilesToMinioDirectly() error {
	// Ensure the bucket exists
	if err := ensureMinIOBucketExists(); err != nil {
		return err
	}

	// Read the test file content
	content, err := os.ReadFile(filepath.Join("..", "assets", "test4.ifc"))
	if err != nil {
		return fmt.Errorf("failed to read test4.ifc file: %w", err)
	}

	// Initialize Kafka writer
	kafkaWriter := newKafkaWriter()
	defer kafkaWriter.Close()

	for _, msg := range messages {
		// Upload file to MinIO with metadata
		contentReader := bytes.NewReader(content)

		// Create metadata similar to the API route
		userMetadata := map[string]string{
			"X-Amz-Meta-Project-Name": msg.Project,
			"X-Amz-Meta-Filename":     msg.FilenameOriginal,
			"X-Amz-Meta-Created-At":   msg.Timestamp,
		}

		objectInfo, err := minioClient.PutObject(
			context.Background(),
			ifcBucket,
			msg.FileNameUUID,
			contentReader,
			int64(len(content)),
			minio.PutObjectOptions{
				ContentType:  "application/octet-stream",
				UserMetadata: userMetadata,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to upload file to MinIO: %w", err)
		}

		log.Printf("Uploaded file to MinIO: %s, size: %d with metadata", objectInfo.Key, objectInfo.Size)

		// Create a simple download link message
		// Format: http://minio:9000/bucket-name/object-path
		minioEndpoint := getEnv("MINIO_HOST", "minio")
		minioPort := getEnv("MINIO_PORT", "9000")
		downloadLink := fmt.Sprintf("http://%s:%s/%s/%s",
			minioEndpoint,
			minioPort,
			ifcBucket,
			msg.FileNameUUID)

		log.Printf("Created download link: %s", downloadLink)

		// Send simple message to Kafka with just the download link
		err = kafkaWriter.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(msg.Project),
				Value: []byte(downloadLink),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to send message to Kafka: %w", err)
		}

		log.Printf("Sent message to Kafka topic '%s': %s", topic, downloadLink)
	}

	return nil
}

func getFileFromPbiServer(baseUrl string, blobID string) ([]byte, error) {
	url := fmt.Sprintf("%s/blob?id=%s", baseUrl, blobID)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add API key to header
	apiKey := getEnv("PBI_API_KEY", "123")
	req.Header.Set("X-API-Key", apiKey)

	// Make the request
	log.Printf("Getting fragments file from PBI server: %s\n", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
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
	log.Println("Adding IFC files to MinIO directly and sending Kafka messages...")
	if err := addIfcFilesToMinioDirectly(); err != nil {
		log.Fatalf("Error adding IFC files to MinIO and sending Kafka messages: %v", err)
	}

	log.Println("Waiting for IFC consumer to convert ifc to gz and write back to MinIO...")
	time.Sleep(5 * time.Minute)

	log.Println("Getting fragments files...")

	for _, msg := range messages {
		fileID := strings.Split(msg.FileNameUUID, ".")[0]
		fileIDgz := fileID + ".gz"

		_, err := getFileFromPbiServer(pbiServerUrl, fileIDgz)
		if err != nil {
			log.Fatalf("Error getting fragments files: %v", err)
		}

		log.Printf("Received fragments file: %s\n", fileIDgz)
	}

	log.Println("IFC blob test completed successfully. Now go and check the updates table in the database to see if the data about the blob has been written.")
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
