package azure

import (
	"viz_pbi-server/env"
)

// Config holds Azure Blob Storage configuration
type Config struct {
	ConnectionString string
	Container        string
}

// NewConfig creates a new Azure Blob Storage configuration from environment variables
func NewConfig() *Config {
	// Load required environment variables
	connectionString := env.Get("AZURE_STORAGE_CONNECTION_STRING")
	container := env.Get("VIZ_IFC_FRAGMENTS_BUCKET")

	return &Config{
		ConnectionString: connectionString,
		Container:        container,
	}
}
