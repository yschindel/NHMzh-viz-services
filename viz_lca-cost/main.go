package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"viz_lca-cost/server"
)

func main() {
	// Kafka configuration
	envBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	envTopic := getEnv("KAFKA_LCA_TOPIC", "lca-data")
	costTopic := getEnv("KAFKA_COST_TOPIC", "cost-data")
	groupID := getEnv("VIZ_KAFKA_DATA_GROUP_ID", "lca-consumer-group")

	// Azure configuration
	azureServer := getEnv("AZURE_DB_SERVER", "")
	azurePort := getEnv("AZURE_DB_PORT", "")
	azureUser := getEnv("AZURE_DB_USER", "")
	azurePassword := getEnv("AZURE_DB_PASSWORD", "")
	azureDatabase := getEnv("AZURE_DB_DATABASE", "")

	config := server.DBConfig{
		Server:   azureServer,
		Port:     StringToInt(azurePort),
		User:     azureUser,
		Password: azurePassword,
		Database: azureDatabase,
	}

	azureDB, err := server.ConnectDB(config)
	if err != nil {
		log.Fatal(err)
	}
	defer azureDB.Close()

	// Create and start consumer
	consumer := server.NewConsumer(
		envBroker,
		envTopic,
		envBroker, // Using same broker for both topics
		costTopic,
		groupID,
		azureDB,
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

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("failed to convert string to int: %v", err)
	}
	return i
}
