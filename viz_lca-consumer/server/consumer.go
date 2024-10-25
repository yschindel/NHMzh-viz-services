package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	writer *Writer
}

type LcaMessage struct {
	Project   string    `json:"project"`
	Filename  string    `json:"filename"`
	Timestamp string    `json:"timestamp"`
	Data      []Element `json:"data"`
}

type Element struct {
	Category   string  `json:"category"`
	CO2e       float64 `json:"co2e"`
	GreyEnergy float64 `json:"greyEnergy"`
	UBP        float64 `json:"UBP"`
}

func NewConsumer(broker, topic, groupID string, writer *Writer) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})

	log.Printf("consumer created for broker: %s, topic: %s, groupID: %s", broker, topic, groupID)

	return &Consumer{
		reader: r,
		writer: writer,
	}
}

func (c *Consumer) StartConsuming(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("could not read message: %v", err)
			continue
		}
		c.handleMessage(m)
	}
}

func (c *Consumer) handleMessage(m kafka.Message) {
	log.Printf("received message: %s", string(m.Key))

	var message LcaMessage
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		log.Printf("could not unmarshal message: %v", err)
		return
	}

	log.Printf("Project: %v \n", message.Project)
	log.Printf("Filename: %v \n", message.Filename)
	log.Printf("Timestamp: %v \n", message.Timestamp)
	log.Printf("Data Count: %v \n", len(message.Data))

	_, err = c.writer.InsertDoc(message)
	if err != nil {
		log.Printf("could not insert document: %v", err)
	}
}
