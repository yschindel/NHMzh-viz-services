package main

import (
	"encoding/json"
	"fmt"
	"lca_cost/data"
	"lca_cost/mongo"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

var kafkaBroker = "localhost:9092"
var lcaTopic string
var costTopic string
var pbiServerURL string

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

	pbiServerPort := os.Getenv("PBI_SERVER_PORT")
	if pbiServerPort == "" {
		pbiServerPort = "8080"
	}
	pbiServerURL = "http://localhost:" + pbiServerPort
}

func getFromPbiServer(url string, db string, collection string) ([]mongo.Element, error) {
	// Make HTTP request to /data endpoint
	resp, err := http.Get(url + "/data?db=" + db + "&collection=" + collection)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to PBI server: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from PBI server: got %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// Decode the response body into a slice of Elements
	var elements []mongo.Element
	if err := json.NewDecoder(resp.Body).Decode(&elements); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return elements, nil
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

	time.Sleep(5 * time.Second)

	for msgIndex, msg := range data.LcaMessages {
		log.Printf("checking data for %s/%s at %s", msg.Project, msg.Filename, msg.Timestamp)
		dbEls, err := getFromPbiServer(pbiServerURL, msg.Project, msg.Filename)
		if err != nil {
			log.Fatalf("Error getting LCA data from PBI server: %v", err)
		}

		// filter data for the timestamp
		dbEls = func() []mongo.Element {
			var filteredEls []mongo.Element
			for _, el := range dbEls {
				if el.Timestamp == msg.Timestamp {
					filteredEls = append(filteredEls, el)
				}
			}
			return filteredEls
		}()

		if len(dbEls) != len(msg.Data) {
			log.Fatalf("Number of elements mismatch: %d != %d", len(dbEls), len(msg.Data))
		}

		dbElsMap := make(map[string]mongo.Element)
		for _, dbEl := range dbEls {
			dbElsMap[dbEl.Id] = dbEl
		}

		for dataIndex, msg := range msg.Data {
			dbEl, ok := dbElsMap[msg.Id]
			if !ok {
				log.Printf("Element not found in DB: %v", msg)
				continue
			}

			if !data.IsLcaItemEqualTo(msg, dbEl) {
				log.Printf("Element Lca mismatch: %v != %v", msg, dbEl)
				continue
			}

			if !data.IsCostItemEqualTo(data.CostMessages[msgIndex].Data[dataIndex], dbEl) {
				log.Printf("Element cost mismatch: %v != %v", data.CostMessages[msgIndex].Data[dataIndex], dbEl)
				continue
			}
		}

		log.Printf("Data for %s/%s at %s match", msg.Project, msg.Filename, msg.Timestamp)
		log.Println("viz_pbi-server: Test passed!")
	}

	// Wait for all goroutines to finish
	<-done
	<-done
}
