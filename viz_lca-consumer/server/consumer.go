package server

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	writer *Writer
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
	log.Printf("received message: %s", string(m.Value))

	document := map[string]interface{}{
		"message":   string(m.Value),
		"topic":     m.Topic,
		"partition": m.Partition,
		"offset":    m.Offset,
	}

	_, err := c.writer.InsertDocument(document)
	if err != nil {
		log.Printf("could not insert document: %v", err)
	}
}
