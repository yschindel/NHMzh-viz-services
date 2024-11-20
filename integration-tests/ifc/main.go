package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaBroker     = "localhost:9092" // has to be hardcoded for the test because the .env uses 'docker-compose' service name
	topic           string
	ifcBucket       string
	fragmentsBucket string
	minioEndpoint   = "localhost" // has to be hardcoded for the test because the .env uses 'docker-compose' service name
	minioPort       int
	minioUrl        string
	minioUseSSL     bool
	minioAccessKey  string
	minioSecretKey  string
	pbiServerPort   int
	pbiServerUrl    string
	messages        []Message
	minioClient     *minio.Client
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
	minioPort = getEnvAsInt("MINIO_PORT", 9000)
	minioUseSSL = getEnvAsBool("MINIO_USE_SSL", false)
	minioAccessKey = getEnv("MINIO_ACCESS_KEY", "")
	minioSecretKey = getEnv("MINIO_SECRET_KEY", "")
	pbiServerPort = getEnvAsInt("PBI_SERVER_PORT", 3000)
	pbiServerUrl = "http://localhost:" + strconv.Itoa(pbiServerPort)
	minioUrl = fmt.Sprintf("%s:%d", minioEndpoint, minioPort)

	// log all environment variables
	log.Printf("KAFKA_BROKER: %s", kafkaBroker)
	log.Printf("KAFKA_IFC_TOPIC: %s", topic)
	log.Printf("MINIO_IFC_BUCKET: %s", ifcBucket)
	log.Printf("MINIO_FRAGMENTS_BUCKET: %s", fragmentsBucket)
	log.Printf("MINIO_ENDPOINT: %s", minioEndpoint)
	log.Printf("MINIO_PORT: %d", minioPort)
	log.Printf("MINIO_USE_SSL: %t", minioUseSSL)
	log.Printf("MINIO_ACCESS_KEY: %s", minioAccessKey)
	log.Printf("MINIO_SECRET_KEY: %s", minioSecretKey)
	log.Printf("PBI_SERVER_PORT: %d", pbiServerPort)
	log.Printf("MINIO_URL: %s", minioUrl)

	// Initialize messages
	messages = []Message{
		newMessage("project1", "file1.ifc"),
		newMessage("project2", "file2.ifc"),
		newMessage("project1", "file3.ifc"),
	}

	client, err := newMinioClient()
	if err != nil {
		log.Fatalf("Error creating MinIO client: %v", err)
	}
	minioClient = client

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

func newMinioClient() (*minio.Client, error) {
	client, err := minio.New(minioUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: minioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}
	return client, nil
}

func newFilename(project, filename, timestamp string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	fileTimestamp := timestamp
	return fmt.Sprintf("%s/%s_%s%s", project, name, fileTimestamp, ext)
}

func addIfcFileToMinio(client *minio.Client, location string, content []byte) error {
	_, err := client.PutObject(context.Background(), ifcBucket, location, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
	return err
}

func addIfcFilesToMinio() error {
	if err := createBucket(minioClient, ifcBucket); err != nil {
		return err
	}

	content, err := os.ReadFile(filepath.Join("..", "assets", "test.ifc"))
	if err != nil {
		return fmt.Errorf("failed to read test.ifc file: %w", err)
	}

	for _, msg := range messages {
		if err := addIfcFileToMinio(minioClient, msg.Location, content); err != nil {
			return fmt.Errorf("failed to add file to MinIO: %w", err)
		}
		log.Printf("Added file to MinIO: %s\n", msg.Location)
	}

	return nil
}

func createBucket(client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	}
	return nil
}

func produceMessages() error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	for _, msg := range messages {
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Value: data,
		})
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		log.Printf("Sent message: %s/%s\n", msg.Project, msg.Filename)
	}

	return nil
}

func consumeMessages() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		GroupID:  "test-group",
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var value Message
		if err := json.Unmarshal(msg.Value, &value); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		log.Printf("Received message: %s/%s\n", value.Project, value.Filename)
	}

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

	log.Println("Waiting for MinIO to process IFC file...")
	time.Sleep(2 * time.Second)

	log.Println("Producing IFC messages...")
	if err := produceMessages(); err != nil {
		log.Fatalf("Error producing IFC messages: %v", err)
	}

	log.Println("Consuming IFC messages...")
	go func() {
		if err := consumeMessages(); err != nil {
			log.Fatalf("Error consuming IFC messages: %v", err)
		}
	}()

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

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
