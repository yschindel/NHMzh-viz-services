package env

import (
	"os"

	"viz_lca-cost/logger"
)

// Get retrieves an environment variable or returns an empty string if not found
func Get(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		logger.Error("Environment variable not set", "env_var", key)
		// exit without restarting the container
		os.Exit(0)
	} else {
		logger.Debug("Environment variable loaded", "env_var", key)
	}
	return value
}
