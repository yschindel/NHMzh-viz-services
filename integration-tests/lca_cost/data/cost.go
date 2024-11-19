package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var CostMessages []CostMessage

// CostMessage represents the structure of an cost message
type CostMessage struct {
	Project   string         `json:"project"`
	Filename  string         `json:"filename"`
	Data      []CostDataItem `json:"data"`
	Timestamp string         `json:"timestamp"`
}

// costElement represents an element in the cost message
type CostDataItem struct {
	Id       string  `json:"id"`
	Category string  `json:"category"`
	Cost     float32 `json:"cost"`
	CostUnit float32 `json:"cost_unit"`
}

// newMessage creates a new costMessage
func newCostMessage(project, filename string) CostMessage {
	timestamp := time.Now().Format(time.RFC3339)
	return CostMessage{
		Project:   project,
		Filename:  filename,
		Data:      newCostData(DataItems),
		Timestamp: timestamp,
	}
}

// newCostData generates a random list of Elements
func newCostData(dataItems []DataItem) []CostDataItem {
	elements := make([]CostDataItem, len(dataItems))
	for i, item := range dataItems {
		elements[i] = CostDataItem{
			Id:       item.Id,
			Category: item.Category,
			Cost:     randFloat(),
			CostUnit: 200,
		}
	}
	return elements
}

// func IsCostItemEqualTo(a CostDataItem, b mongo.Element) bool {
// 	return a.Id == b.Id && a.Category == b.Category && a.Cost == b.Cost
// }

// produceMessages sends messages to a Kafka topic
func ProduceCostMessages(broker string, topic string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	CostMessages = []CostMessage{
		newCostMessage("project1", "file1.ifc"),
		newCostMessage("project2", "file2.ifc"),
		newCostMessage("project1", "file3.ifc"),
	}

	for _, msg := range CostMessages {
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
func ConsumeCostMessages(broker string, topic string) error {
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
