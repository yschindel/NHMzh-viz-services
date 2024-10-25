package main

import (
	"lca_cost/data"
	"lca_cost/mongo"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

var kafkaBroker = "localhost:9092"
var lcaTopic string
var costTopic string
var mongoReader mongo.Reader

func init() {
	// Load environment variables from .env file
	envPath := filepath.Join("..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	lcaTopic = os.Getenv("KAFKA_LCA_TOPIC")
	if lcaTopic == "" {
		lcaTopic = "lca-data"
	}
	costTopic = os.Getenv("KAFKA_COST_TOPIC")
	if costTopic == "" {
		costTopic = "cost-data"
	}
	mongoReader = *mongo.NewReader("mongodb://localhost:27017")
}

func main() {
	// Channel to signal when both produce and consume are done
	done := make(chan bool)

	// Run produceMessages in a goroutine
	go func() {
		log.Println("Producing LCA messages...")
		if err := data.ProduceLcaMessages(kafkaBroker, lcaTopic); err != nil {
			log.Fatalf("Error producing LCA messages: %v", err)
		}
		done <- true
	}()

	// Run consumeMessages in a goroutine
	go func() {
		log.Println("Consuming LCA messages...")
		if err := data.ConsumeLcaMessages(kafkaBroker, lcaTopic); err != nil {
			log.Fatalf("Error consuming LCA messages: %v", err)
		}
		done <- true
	}()

	// Run produceMessages in a goroutine
	go func() {
		log.Println("Producing Cost messages...")
		if err := data.ProduceCostMessages(kafkaBroker, costTopic); err != nil {
			log.Fatalf("Error producing Cost messages: %v", err)
		}
		done <- true
	}()

	// Run consumeMessages in a goroutine
	go func() {
		log.Println("Consuming Cost messages...")
		if err := data.ConsumeCostMessages(kafkaBroker, costTopic); err != nil {
			log.Fatalf("Error consuming Cost messages: %v", err)
		}
		done <- true
	}()

	time.Sleep(20 * time.Second)

	go func() {
		for _, msg := range data.LcaMessages {
			dbEls, err := mongoReader.ReadAllElements(msg.Project, msg.Filename, msg.Timestamp)
			if err != nil {
				log.Fatalf("Error reading all elements: %v", err)
			}

			dbElsMap := make(map[string]mongo.Element)
			for _, dbEl := range dbEls {
				dbElsMap[dbEl.Id] = dbEl
			}

			for _, el := range msg.Data {
				dbEl, ok := dbElsMap[el.Id]
				if !ok {
					log.Printf("Element not found in DB: %v", el)
					continue
				}

				if !data.IsLcaItemEqualTo(el, dbEl) {
					log.Printf("Element Lca mismatch: %v != %v", el, dbEl)
					continue
				}
			}
			log.Printf("LCA data for %s/%s at %s match", msg.Project, msg.Filename, msg.Timestamp)
		}
	}()

	go func() {
		for _, msg := range data.CostMessages {
			dbEls, err := mongoReader.ReadAllElements(msg.Project, msg.Filename, msg.Timestamp)
			if err != nil {
				log.Fatalf("Error reading all elements: %v", err)
			}

			dbElsMap := make(map[string]mongo.Element)
			for _, dbEl := range dbEls {
				dbElsMap[dbEl.Id] = dbEl
			}

			for _, el := range msg.Data {
				dbEl, ok := dbElsMap[el.Id]
				if !ok {
					log.Printf("Element not found in DB: %v", el)
					continue
				}

				if !data.IsCostItemEqualTo(el, dbEl) {
					log.Printf("Element cost mismatch: %v != %v", el, dbEl)
					continue
				}
			}
			log.Printf("Cost data for %s/%s at %s match", msg.Project, msg.Filename, msg.Timestamp)
		}
	}()

	// Wait for all goroutines to finish
	<-done
	<-done
	<-done
	<-done
}
