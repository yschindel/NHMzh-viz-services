package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	environmentalReader *kafka.Reader
	costReader          *kafka.Reader
	azureDB             *sql.DB
}

func NewConsumer(envBroker, envTopic, costBroker, costTopic, groupID string, azureDB *sql.DB) *Consumer {
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
		azureDB:             azureDB,
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

	message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	err = WriteLcaMessage(c.azureDB, message)
	if err != nil {
		log.Printf("could not write lca message: %v", err)
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

	message.Filename = strings.TrimSuffix(message.Filename, ".ifc")

	err = WriteCostMessage(c.azureDB, message)
	if err != nil {
		log.Printf("could not write cost message: %v", err)
		return
	}
}
