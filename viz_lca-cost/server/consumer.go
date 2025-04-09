package server

import (
	"context"
	"encoding/json"
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
	// Start both consumers in separate goroutines
	go c.consumeLca(ctx)
	go c.consumeCost(ctx)

	// Keep the main goroutine alive
	<-ctx.Done()
}

func (c *Consumer) consumeLca(ctx context.Context) {
	log := c.logger.WithFields(logger.Fields{
		"consumer": "lca",
		"topic":    c.lcaReader.Config().Topic,
	})
	log.Info("Starting consumer")

	for {
		m, err := c.lcaReader.ReadMessage(ctx)
		if err != nil {
			log.Error("Could not read message", logger.Fields{
				"error": err,
			})
			continue
		}
		log.Info("Received message", logger.Fields{
			"key":    string(m.Key),
			"length": len(m.Value),
		})
		c.handleEnvironmentalMessage(m)
	}
}

func (c *Consumer) consumeCost(ctx context.Context) {
	log := c.logger.WithFields(logger.Fields{
		"consumer": "cost",
		"topic":    c.costReader.Config().Topic,
	})
	log.Info("Starting consumer")

	for {
		m, err := c.costReader.ReadMessage(ctx)
		if err != nil {
			log.Error("Could not read message", logger.Fields{
				"error": err,
			})
			continue
		}
		log.Info("Received message", logger.Fields{
			"key":    string(m.Key),
			"length": len(m.Value),
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
	// This could be separated into another service
	// In that case we would POST to an API endpoint here.
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
