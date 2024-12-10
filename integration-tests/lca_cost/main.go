package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lca_cost/data"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/marcboeker/go-duckdb"
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

func getModelListFromPbiServer(url string, project string) ([]string, error) {
	resp, err := http.Get(url + "/data/models?project=" + project)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to PBI server: %v", err)
	}
	defer resp.Body.Close()

	var modelList []string
	if err := json.NewDecoder(resp.Body).Decode(&modelList); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return modelList, nil
}

func getProjectListFromPbiServer(url string) ([]string, error) {
	resp, err := http.Get(url + "/data/projects")
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to PBI server: %v", err)
	}
	defer resp.Body.Close()

	var projectList []string
	if err := json.NewDecoder(resp.Body).Decode(&projectList); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return projectList, nil
}

// get a parquet file from the lca-cost-data bucket
func getDataFromPbiServer(url string, project string, filename string) ([]byte, error) {
	resp, err := http.Get(url + "/data?name=" + filename + "&project=" + project)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to PBI server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read all bytes at once, similar to arrayBuffer() in TypeScript
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return data, nil
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

	// time.Sleep(5 * time.Second)

	// for _, msg := range data.LcaMessages {
	// 	log.Printf("checking data for %s/%s at %s", msg.Project, msg.Filename, msg.Timestamp)
	// 	filename := strings.Replace(msg.Filename, ".ifc", ".parquet", 1)

	// 	log.Printf("getting parquet file from PBI server: %s", filename)
	// 	parquetBytes, err := getDataFromPbiServer(pbiServerURL, msg.Project, filename)
	// 	if err != nil {
	// 		log.Fatalf("Error getting LCA data from PBI server: %v", err)
	// 	}

	// 	db, err := sql.Open("duckdb", "")
	// 	if err != nil {
	// 		log.Fatalf("failed to open duckdb: %v", err)
	// 	}

	// 	// write the parquet file to disk not a temp file
	// 	parquetFile, err := os.Create(filename)
	// 	if err != nil {
	// 		log.Fatalf("failed to create parquet file: %v", err)
	// 	}
	// 	defer parquetFile.Close()
	// 	parquetFile.Write(parquetBytes)

	// 	rows, err := db.Query(`SELECT * FROM read_parquet($1) LIMIT 3`, filename)
	// 	if err != nil {
	// 		log.Fatalf("failed to query parquet file: %v", err)
	// 	}
	// 	defer rows.Close()

	// 	// Get column names
	// 	columns, err := rows.Columns()
	// 	if err != nil {
	// 		log.Fatalf("failed to get columns: %v", err)
	// 	}
	// 	log.Printf("Columns: %v", columns)

	// 	// Print each row
	// 	for rows.Next() {
	// 		// Create a slice of interface{} to hold the values
	// 		values := make([]interface{}, len(columns))
	// 		valuePtrs := make([]interface{}, len(columns))

	// 		// Create pointers to each interface{}
	// 		for i := range values {
	// 			valuePtrs[i] = &values[i]
	// 		}

	// 		// Scan the result into the pointers
	// 		if err := rows.Scan(valuePtrs...); err != nil {
	// 			log.Fatalf("failed to scan row: %v", err)
	// 		}

	// 		// Print the row data
	// 		rowData := make(map[string]interface{})
	// 		for i, col := range columns {
	// 			val := values[i]
	// 			rowData[col] = val
	// 		}
	// 		log.Printf("Row: %+v", rowData)
	// 	}

	// 	if err = rows.Err(); err != nil {
	// 		log.Fatalf("error iterating rows: %v", err)
	// 	}

	// 	log.Printf("Data for %s/%s at %s match", msg.Project, msg.Filename, msg.Timestamp)
	// 	log.Println("viz_pbi-server: Test passed!")
	// }

	// modelList, err := getModelListFromPbiServer(pbiServerURL, "project1")
	// if err != nil {
	// 	log.Fatalf("Error getting model list from PBI server: %v", err)
	// }

	// log.Printf("modelList: %v", modelList)

	// projectList, err := getProjectListFromPbiServer(pbiServerURL)
	// if err != nil {
	// 	log.Fatalf("Error getting project list from PBI server: %v", err)
	// }

	// log.Printf("projectList: %v", projectList)

	// Wait for all goroutines to finish
	<-done
	<-done
}
