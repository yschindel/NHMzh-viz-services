package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"viz_lca-cost/env"

	"github.com/rs/zerolog/log"
)

type MessageWriter struct {
	url          string
	endpointLca  string
	endpointCost string
	apiKey       string
}

func NewMessageWriter() *MessageWriter {
	url := env.Get("PBI_SERVER_URL")
	apiKey := env.Get("PBI_SRV_API_KEY")
	epLca := url + env.Get("STORAGE_DATA_LCA_ENDPOINT")
	epCost := url + env.Get("STORAGE_DATA_COST_ENDPOINT")

	log.Info().Msgf("LCA endpoint: %s", epLca)
	log.Info().Msgf("Cost endpoint: %s", epCost)

	return &MessageWriter{url: url, endpointLca: epLca, endpointCost: epCost, apiKey: apiKey}
}

// WriteLcaMessage writes a message to the LCA endpoint
func (w *MessageWriter) WriteLcaMessage(message []EavMaterialDataItem) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return w.WriteMessage(jsonData, w.endpointLca)
}

// WriteCostMessage writes a message to the cost endpoint
func (w *MessageWriter) WriteCostMessage(message []EavElementDataItem) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return w.WriteMessage(jsonData, w.endpointCost)
}

// WriteMessage writes a message to a given url
func (w *MessageWriter) WriteMessage(jsonData []byte, url string) error {
	// Create request with API key header
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set the content type and API key header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", w.apiKey)

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
