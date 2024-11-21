package server

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	environmentalReader *kafka.Reader
	costReader          *kafka.Reader
	db                  *DuckDBManager
}

func NewConsumer(envBroker, envTopic, costBroker, costTopic, groupID string, db *DuckDBManager) *Consumer {
	environmentalReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{envBroker},
		Topic:   envTopic,
		GroupID: groupID + "-env",
	})

	costReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{costBroker},
		Topic:   costTopic,
		GroupID: groupID + "-cost",
	})

	log.Printf("environmental consumer created for topic: %s", envTopic)
	log.Printf("cost consumer created for topic: %s", costTopic)

	return &Consumer{
		environmentalReader: environmentalReader,
		costReader:          costReader,
		db:                  db,
	}
}

func (c *Consumer) StartConsuming(ctx context.Context) {
	// Start both consumers in separate goroutines
	go c.consumeEnvironmental(ctx)
	go c.consumeCost(ctx)

	// Keep the main goroutine alive
	<-ctx.Done()
}

func (c *Consumer) consumeEnvironmental(ctx context.Context) {
	log.Printf("Starting environmental consumer...")
	for {
		m, err := c.environmentalReader.ReadMessage(ctx)
		if err != nil {
			log.Printf("could not read environmental message: %v", err)
			continue
		}
		log.Printf("Received environmental message with key: %s, length: %d bytes", string(m.Key), len(m.Value))
		c.handleEnvironmentalMessage(m)
	}
}

func (c *Consumer) consumeCost(ctx context.Context) {
	log.Printf("Starting cost consumer...")
	for {
		m, err := c.costReader.ReadMessage(ctx)
		if err != nil {
			log.Printf("could not read cost message: %v", err)
			continue
		}
		log.Printf("Received cost message with key: %s, length: %d bytes", string(m.Key), len(m.Value))
		c.handleCostMessage(m)
	}
}

func (c *Consumer) handleEnvironmentalMessage(m kafka.Message) {
	log.Printf("received environmental message: %s", string(m.Key))

	var message LcaMessage
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		log.Printf("could not unmarshal environmental message: %v", err)
		return
	}

	filename := strings.TrimSuffix(message.Filename, ".ifc")

	// Ensure parquet file exists
	err = c.db.EnsureParquetFile(message.Project, filename)
	if err != nil {
		log.Printf("could not ensure parquet file: %v", err)
		return
	}

	// Update environmental data with timestamp
	err = c.db.UpdateParquetFromEnvironmentalData(message.Project, filename, message.Data, message.Timestamp)
	if err != nil {
		log.Printf("could not update environmental data: %v", err)
		return
	}
}

func (c *Consumer) handleCostMessage(m kafka.Message) {
	log.Printf("received cost message: %s", string(m.Key))

	var message CostMessage
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		log.Printf("could not unmarshal cost message: %v", err)
		return
	}

	filename := strings.TrimSuffix(message.Filename, ".ifc")

	// Ensure parquet file exists
	err = c.db.EnsureParquetFile(message.Project, filename)
	if err != nil {
		log.Printf("could not ensure parquet file: %v", err)
		return
	}

	// Update cost data with timestamp
	err = c.db.UpdateParquetFromCostData(message.Project, filename, message.Data, message.Timestamp)
	if err != nil {
		log.Printf("could not update cost data: %v", err)
		return
	}
}
