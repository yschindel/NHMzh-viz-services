package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

var LcaMessages []LcaMessage

// Message represents the structure of an LCA message
type LcaMessage struct {
	Project   string        `json:"project"`
	Filename  string        `json:"filename"`
	Data      []LcaDataItem `json:"data"`
	Timestamp string        `json:"timestamp"`
}

// LcaDataItem represents an element in the LCA message
type LcaDataItem struct {
	Id           string  `json:"id"`
	Category     string  `json:"ebkph"`
	MaterialKbob string  `json:"mat_kbob"`
	GwpAbsolute  float32 `json:"gwp_absolute"`
	GwpRelative  float32 `json:"gwp_relative"`
	PenrAbsolute float32 `json:"penr_absolute"`
	PenrRelative float32 `json:"penr_relative"`
	UbpAbsolute  float32 `json:"ubp_absolute"`
	UbpRelative  float32 `json:"ubp_relative"`
}

var MaterialKbob = [6]string{"Aether", "Coaxium", "Dilithium", "Kryptonite", "Spice", "Vibranium"}

// newMessage creates a new LcaMessage
func newLcaMessage(project, filename string) LcaMessage {
	timestamp := time.Now().Format(time.RFC3339)
	return LcaMessage{
		Project:   project,
		Filename:  filename,
		Data:      newLcaData(DataItems),
		Timestamp: timestamp,
	}
}

// newData generates a list of Elements based on provided IDs
func newLcaData(dataItems []DataItem) []LcaDataItem {
	elements := make([]LcaDataItem, len(dataItems))
	for i, item := range dataItems {
		elements[i] = LcaDataItem{
			Id:           item.Id,
			Category:     item.Category,
			MaterialKbob: MaterialKbob[rand.Intn(len(MaterialKbob))],
			GwpAbsolute:  randFloat(),
			GwpRelative:  randFloat(),
			PenrAbsolute: randFloat(),
			PenrRelative: randFloat(),
			UbpAbsolute:  randFloat(),
			UbpRelative:  randFloat(),
		}
	}
	return elements
}

// func IsLcaItemEqualTo(a LcaDataItem, b mongo.Element) bool {
// 	return a.Id == b.Id && a.Category == b.Category && a.CO2e == b.CO2e && a.GreyEnergy == b.GreyEnergy && a.UBP == b.UBP
// }

// produceMessages sends messages to a Kafka topic
func ProduceLcaMessages(broker string, topic string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	LcaMessages = []LcaMessage{
		newLcaMessage("project1", "file1.ifc"),
		newLcaMessage("project2", "file2.ifc"),
		newLcaMessage("project1", "file3.ifc"),
	}

	for _, msg := range LcaMessages {
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

// consumeMessages reads messages from a Kafka topic
func ConsumeLcaMessages(broker string, topic string) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
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

		var value CostMessage
		if err := json.Unmarshal(msg.Value, &value); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		log.Printf("Received message: %s/%s\n", value.Project, value.Filename)
	}
}
