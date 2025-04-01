package azure

import (
	"os"
)

type Config struct {
	AccountName string
	AccountKey  string
	EndpointURL string
}

func NewConfig() *Config {
	return &Config{
		AccountName: getEnv("AZURE_STORAGE_ACCOUNT", "devstoreaccount1"),
		AccountKey:  getEnv("AZURE_STORAGE_KEY", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="),
		EndpointURL: getEnv("AZURE_STORAGE_URL", "http://127.0.0.1:10000/devstoreaccount1"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
