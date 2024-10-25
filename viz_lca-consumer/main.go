package main

import (
	"context"
	"os"
	"viz_lca-consumer/server"
)

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_LCA_TOPIC", "your-topic")
	groupID := getEnv("VIZ_KAFKA_LCA_GROUP_ID", "your-group-id")

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")

	writer := server.NewWriter(mongoURI)
	consumer := server.NewConsumer(broker, topic, groupID, writer)

	ctx := context.Background()
	consumer.StartConsuming(ctx)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
