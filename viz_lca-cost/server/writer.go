package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"viz_lca-cost/env"
	"viz_lca-cost/logger"
)

type MessageWriter struct {
	url          string
	endpointLca  string
	endpointCost string
	apiKey       string
	logger       *logger.Logger
}

func NewMessageWriter() *MessageWriter {
	url := env.Get("STORAGE_SERVICE_URL")
	apiKey := env.Get("STORAGE_SERVICE_API_KEY")
	epLca := url + env.Get("STORAGE_DATA_LCA_ENDPOINT")
	epCost := url + env.Get("STORAGE_DATA_COST_ENDPOINT")

	log := logger.WithFields(logger.Fields{
		"lca_endpoint":  epLca,
		"cost_endpoint": epCost,
	})
	log.Info("Initializing message writer")

	return &MessageWriter{
		url:          url,
		endpointLca:  epLca,
		endpointCost: epCost,
		apiKey:       apiKey,
		logger:       log,
	}
}

// WriteLcaMessage writes a message to the LCA endpoint
func (w *MessageWriter) WriteLcaMessage(message []EavMaterialDataItem) error {
	log := w.logger.WithFields(logger.Fields{
		"count": len(message),
	})

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Error("Error marshaling LCA message", logger.Fields{
			"error": err,
		})
		return err
	}

	return w.WriteMessage(jsonData, w.endpointLca)
}

// WriteCostMessage writes a message to the cost endpoint
func (w *MessageWriter) WriteCostMessage(message []EavElementDataItem) error {
	log := w.logger.WithFields(logger.Fields{
		"count": len(message),
	})

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Error("Error marshaling cost message", logger.Fields{
			"error": err,
		})
		return err
	}

	return w.WriteMessage(jsonData, w.endpointCost)
}

// WriteMessage writes a message to a given url
func (w *MessageWriter) WriteMessage(jsonData []byte, url string) error {
	log := w.logger.WithFields(logger.Fields{
		"url": url,
	})

	// Create request with API key header
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("Error creating request", logger.Fields{
			"error": err,
		})
		return err
	}

	// Set the content type and API key header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", w.apiKey)

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending request", logger.Fields{
			"error": err,
		})
		return err
	}
	defer resp.Body.Close()

	log.Info("Successfully sent message")
	return nil
}
