package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	lcaReader  *kafka.Reader
	costReader *kafka.Reader
	processor  *MessageProcessor
	writer     *MessageWriter
}

func NewConsumer(envBroker, envTopic, costBroker, costTopic, groupID string, azureDB *sql.DB) *Consumer {
	lcaReader := kafka.NewReader(kafka.ReaderConfig{
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
		lcaReader:  lcaReader,
		costReader: costReader,
		processor:  NewMessageProcessor(),
		writer:     NewMessageWriter(azureDB),
	}
}

func (c *Consumer) StartConsuming(ctx context.Context) {
	// Start both consumers in separate goroutines
	go c.consumeLca(ctx)
	go c.consumeCost(ctx)

	// Keep the main goroutine alive
	<-ctx.Done()
}

func (c *Consumer) consumeLca(ctx context.Context) {
	log.Printf("Starting environmental consumer...")
	for {
		m, err := c.lcaReader.ReadMessage(ctx)
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

	// Step 1: Receive - already done via Kafka consumer

	// Step 2: Process message
	var message LcaMessage
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		log.Printf("could not unmarshal environmental message: %v", err)
		return
	}

	// Process the message (apply transformations, validations, etc.)
	eavItems, err := c.processor.ProcessLcaMessage(&message)
	if err != nil {
		log.Printf("could not process lca message: %v", err)
		return
	}

	// Step 3: Write to database
	// This could be separated into another service
	// In that case we would POST to an API endpoint here.
	err = c.writer.WriteLcaMessage(eavItems)
	if err != nil {
		log.Printf("could not write lca message: %v", err)
		return
	}
}

func (c *Consumer) handleCostMessage(m kafka.Message) {
	log.Printf("received cost message: %s", string(m.Key))

	// Step 1: Receive - already done via Kafka consumer

	// Step 2: Process message
	var message CostMessage
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		log.Printf("could not unmarshal cost message: %v", err)
		return
	}

	// Process the message (apply transformations, validations, etc.)
	eavItems, err := c.processor.ProcessCostMessage(&message)
	if err != nil {
		log.Printf("could not process cost message: %v", err)
		return
	}

	// Step 3: Write to database
	err = c.writer.WriteCostMessage(eavItems)
	if err != nil {
		log.Printf("could not write cost message: %v", err)
		return
	}
}
