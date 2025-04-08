package main

import (
	"context"
	"viz_lca-cost/env"
	"viz_lca-cost/logger"
	"viz_lca-cost/server"
)

func main() {
	log := logger.New()

	// Kafka configuration
	envBroker := env.Get("KAFKA_BROKER")
	envTopic := env.Get("KAFKA_LCA_TOPIC")
	costTopic := env.Get("KAFKA_COST_TOPIC")
	groupID := "viz_data_consumers_group"

	log.Info("Starting LCA-Cost service", logger.Fields{
		"kafka_broker": envBroker,
		"lca_topic":    envTopic,
		"cost_topic":   costTopic,
		"group_id":     groupID,
	})

	// Create and start consumer
	consumer := server.NewConsumer(
		envBroker,
		envTopic,
		envBroker, // Using same broker for both topics
		costTopic,
		groupID,
	)

	ctx := context.Background()
	consumer.StartConsuming(ctx)
}
