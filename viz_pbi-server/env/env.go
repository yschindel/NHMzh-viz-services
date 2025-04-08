// Simple utility to load a single environment variable
package env

import (
	"os"

	"viz_pbi-server/logger"
)

// Get retrieves an environment variable or returns an empty string if not found
func Get(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		logger.Error("Environment variable not set:", "env_var", key)
		os.Exit(0)
	} else {
		logger.Debug("Environment variable loaded:", "env_var", key)
	}
	return value
}
