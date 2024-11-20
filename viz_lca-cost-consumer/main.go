package main

import (
	"context"
	"log"
	"os"
	"viz_lca-cost-consumer/server"
)

func main() {
	// Kafka configuration
	envBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	envTopic := getEnv("KAFKA_LCA_TOPIC", "lca-data")
	costTopic := getEnv("KAFKA_COST_TOPIC", "cost-data")
	groupID := getEnv("VIZ_KAFKA_DATA_GROUP_ID", "lca-consumer-group")

	// MinIO configuration
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost")
	minioPort := getEnv("MINIO_PORT", "9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "ROOTUSER")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "CHANGEME123")
	minioBucket := getEnv("MINIO_LCA_COST_DATA_BUCKET", "lca-cost-data")

	// Initialize DuckDB manager
	credentials := server.CustomMinioCredentials{
		Endpoint:        minioEndpoint + ":" + minioPort,
		AccessKeyID:     minioAccessKey,
		SecretAccessKey: minioSecretKey,
	}

	db, err := server.NewDuckDBManager("", credentials, minioBucket)
	if err != nil {
		log.Fatalf("failed to create DuckDB manager: %v", err)
	}
	defer db.Close()

	// Create and start consumer
	consumer := server.NewConsumer(
		envBroker,
		envTopic,
		envBroker, // Using same broker for both topics
		costTopic,
		groupID,
		db,
	)

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
