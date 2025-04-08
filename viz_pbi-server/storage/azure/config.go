package azure

import (
	"viz_pbi-server/env"
)

// Config holds Azure Blob Storage configuration
type Config struct {
	AccountName string
	AccountKey  string
	EndpointURL string
	Container   string
}

// NewConfig creates a new Azure Blob Storage configuration from environment variables
func NewConfig() *Config {
	// Load required environment variables
	accountName := env.Get("AZURE_STORAGE_ACCOUNT")
	accountKey := env.Get("AZURE_STORAGE_KEY")
	endpointURL := env.Get("AZURE_STORAGE_URL")
	container := env.Get("VIZ_IFC_FRAGMENTS_BUCKET")

	return &Config{
		AccountName: accountName,
		AccountKey:  accountKey,
		EndpointURL: endpointURL,
		Container:   container,
	}
}
