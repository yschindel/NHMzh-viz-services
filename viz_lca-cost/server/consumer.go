package server

import (
	"context"
	"encoding/json"
	"time"
	"viz_lca-cost/logger"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	lcaReader  *kafka.Reader
	costReader *kafka.Reader
	processor  *MessageProcessor
	writer     *MessageWriter
	logger     *logger.Logger
}

func NewConsumer(envBroker, envTopic, costBroker, costTopic, groupID string) *Consumer {
	log := logger.New()

	lcaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{envBroker},
		Topic:       envTopic,
		GroupID:     groupID + "-env",
		StartOffset: kafka.LastOffset,
	})

	costReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{costBroker},
		Topic:       costTopic,
		GroupID:     groupID + "-cost",
		StartOffset: kafka.LastOffset,
	})

	log.Info("Created consumers", logger.Fields{
		"env_topic":  envTopic,
		"cost_topic": costTopic,
	})

	return &Consumer{
		lcaReader:  lcaReader,
		costReader: costReader,
		processor:  NewMessageProcessor(),
		writer:     NewMessageWriter(),
		logger:     log,
	}
}

func (c *Consumer) StartConsuming(ctx context.Context) {
	// Create error channels for both consumers
	lcaErrChan := make(chan error, 1)
	costErrChan := make(chan error, 1)

	// Start both consumers in separate goroutines
	go c.consumeLca(ctx, lcaErrChan)
	go c.consumeCost(ctx, costErrChan)

	// Wait for either context cancellation or consumer errors
	select {
	case <-ctx.Done():
		c.logger.Warn("Context cancelled, shutting down consumers")
	case err := <-lcaErrChan:
		c.logger.Error("LCA consumer failed", logger.Fields{"error": err})
	case err := <-costErrChan:
		c.logger.Error("Cost consumer failed", logger.Fields{"error": err})
	}

	// Close readers
	if err := c.lcaReader.Close(); err != nil {
		c.logger.Error("Error closing LCA reader", logger.Fields{"error": err})
	}
	if err := c.costReader.Close(); err != nil {
		c.logger.Error("Error closing Cost reader", logger.Fields{"error": err})
	}
}

func (c *Consumer) consumeLca(ctx context.Context, errChan chan<- error) {
	log := c.logger.WithFields(logger.Fields{
		"consumer": "lca",
		"topic":    c.lcaReader.Config().Topic,
	})
	log.Info("Starting consumer")

	// Try to fetch initial message with retries
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		_, err := c.lcaReader.FetchMessage(ctx)
		if err == nil {
			log.Info("Fetched initial message successfully")
			break
		}

		if i == maxRetries-1 {
			log.Error("Could not fetch initial message after retries", logger.Fields{
				"error":   err,
				"retries": i + 1,
			})
			errChan <- err
			return
		}

		log.Warn("Retrying initial message fetch", logger.Fields{
			"retry": i + 1,
			"error": err,
		})

		time.Sleep(3 * time.Second)
	}

	for {
		m, err := c.lcaReader.ReadMessage(ctx)
		if err != nil {
			log.Error("Could not read message", logger.Fields{
				"error": err,
			})
			continue
		}
		log.Info("Received message", logger.Fields{
			"key":       string(m.Key),
			"length":    len(m.Value),
			"offset":    m.Offset,
			"partition": m.Partition,
		})
		c.handleEnvironmentalMessage(m)
	}
}

func (c *Consumer) consumeCost(ctx context.Context, errChan chan<- error) {
	log := c.logger.WithFields(logger.Fields{
		"consumer": "cost",
		"topic":    c.costReader.Config().Topic,
	})
	log.Info("Starting consumer")

	// Try to fetch initial message with retries
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		_, err := c.costReader.FetchMessage(ctx)
		if err == nil {
			log.Info("Fetched initial message successfully")
			break
		}

		if i == maxRetries-1 {
			log.Error("Could not fetch initial message after retries", logger.Fields{
				"error":   err,
				"retries": i + 1,
			})
			errChan <- err
			return
		}

		log.Warn("Retrying initial message fetch", logger.Fields{
			"retry": i + 1,
			"error": err,
		})

		time.Sleep(3 * time.Second)
	}

	for {
		m, err := c.costReader.ReadMessage(ctx)
		if err != nil {
			log.Error("Could not read message", logger.Fields{
				"error": err,
			})
			continue
		}
		log.Info("Received message", logger.Fields{
			"key":       string(m.Key),
			"length":    len(m.Value),
			"offset":    m.Offset,
			"partition": m.Partition,
		})
		c.handleCostMessage(m)
	}
}

func (c *Consumer) handleEnvironmentalMessage(m kafka.Message) {
	log := c.logger.WithFields(logger.Fields{
		"consumer": "lca",
		"key":      string(m.Key),
		"size":     len(m.Value),
	})
	log.Info("Processing message")

	// Step 1: Receive - already done via Kafka consumer

	// Step 2: Process message
	var message LcaMessage
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		log.Error("Could not unmarshal message", logger.Fields{
			"error": err,
		})
		return
	}

	// Process the message (apply transformations, validations, etc.)
	eavItems, err := c.processor.ProcessLcaMessage(&message)
	if err != nil {
		log.Error("Could not process message", logger.Fields{
			"error": err,
		})
		return
	}

	// Step 3: Write to database
	err = c.writer.WriteLcaMessage(eavItems)
	if err != nil {
		log.Error("Could not write message", logger.Fields{
			"error": err,
		})
		return
	}

	log.Info("Successfully processed message")
}

func (c *Consumer) handleCostMessage(m kafka.Message) {
	log := c.logger.WithFields(logger.Fields{
		"consumer": "cost",
		"key":      string(m.Key),
		"size":     len(m.Value),
	})
	log.Info("Processing message")

	// Step 1: Receive - already done via Kafka consumer

	// Step 2: Process message
	var message CostMessage
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		log.Error("Could not unmarshal message", logger.Fields{
			"error": err,
		})
		return
	}

	// Process the message (apply transformations, validations, etc.)
	eavItems, err := c.processor.ProcessCostMessage(&message)
	if err != nil {
		log.Error("Could not process message", logger.Fields{
			"error": err,
		})
		return
	}

	// Step 3: Write to database
	err = c.writer.WriteCostMessage(eavItems)
	if err != nil {
		log.Error("Could not write message", logger.Fields{
			"error": err,
		})
		return
	}

	log.Info("Successfully processed message")
}
