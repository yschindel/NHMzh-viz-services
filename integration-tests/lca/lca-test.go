package main

import (
	"context"
	cryptrand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

// LcaMessage represents the structure of an LCA message
type LcaMessage struct {
	Project   string    `json:"project"`
	Filename  string    `json:"filename"`
	Data      []Element `json:"data"`
	Timestamp string    `json:"timestamp"`
}

// Element represents an element in the LCA message
type Element struct {
	ID         string  `json:"id"`
	Category   string  `json:"category"`
	CO2e       float32 `json:"co2e"`
	GreyEnergy float32 `json:"greyEnergy"`
	UBP        float32 `json:"UBP"`
}

var kafkaBroker = "localhost:9092"
var topic string

func init() {
	// Load environment variables from .env file
	envPath := filepath.Join("..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	topic = os.Getenv("KAFKA_LCA_TOPIC")
	if topic == "" {
		topic = "lca-data"
	}
}

// newMessage creates a new LcaMessage
func newMessage(project, filename string) LcaMessage {
	timestamp := time.Now().Format(time.RFC3339)
	return LcaMessage{
		Project:   project,
		Filename:  filename,
		Data:      newData(),
		Timestamp: timestamp,
	}
}

// newData generates a random list of Elements
func newData() []Element {
	numElements := 9000
	elements := make([]Element, numElements)
	for i := 0; i < numElements; i++ {
		elements[i] = newElement()
	}
	return elements
}

// newElement creates a random Element
func newElement() Element {
	id, err := generate22CharGUID()
	if err != nil {
		log.Fatalf("Error generating GUID: %v", err)
	}
	return Element{
		ID:         id,
		Category:   randCat(),
		CO2e:       randFloat(),
		GreyEnergy: randFloat(),
		UBP:        randFloat(),
	}
}

func generate22CharGUID() (string, error) {
	uuid := make([]byte, 16) // UUID is 16 bytes (128 bits)

	// Generate random bytes
	_, err := cryptrand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version and variant bits for the UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is RFC4122

	// Encode to base64 URL-safe and remove the trailing '=='
	encoded := base64.RawURLEncoding.EncodeToString(uuid)

	return encoded, nil
}

// randFloat generates a random float number between min and max
func randFloat() float32 {
	return rand.Float32()*(100-5) + 5
}

// randCat picks a random category from a predefined list
func randCat() string {
	categories := []string{"IfcWalls", "IfcDoors", "IfcFloors"}
	return categories[rand.Intn(len(categories))]
}

// produceMessages sends messages to a Kafka topic
func produceMessages() error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	messages := []LcaMessage{
		newMessage("project1", "file1.ifc"),
		newMessage("project2", "file2.ifc"),
		newMessage("project1", "file3.ifc"),
	}

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
		fmt.Printf("Sent message: %s/%s\n", msg.Project, msg.Filename)
	}
	return nil
}

// consumeMessages reads messages from a Kafka topic
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

		var value LcaMessage
		if err := json.Unmarshal(msg.Value, &value); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		fmt.Printf("Received message: %s/%s\n", value.Project, value.Filename)
	}
}

func main() {
	// Channel to signal when both produce and consume are done
	done := make(chan bool)

	// Run produceMessages in a goroutine
	go func() {
		fmt.Println("Producing LCA messages...")
		if err := produceMessages(); err != nil {
			log.Fatalf("Error producing messages: %v", err)
		}
		done <- true
	}()

	// Run consumeMessages in a goroutine
	go func() {
		fmt.Println("\nConsuming LCA messages...")
		if err := consumeMessages(); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
		done <- true
	}()

	// Wait for both goroutines to finish
	<-done
	<-done
}
