// Simple utility to load a single environment variable
package env

import (
	"os"

	"viz_pbi-server/logger"
)

// Get loads a single environment variable
func Get(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		logger.Error("Environment variable not set: %s", key)
		os.Exit(1)
	} else {
		logger.Debug("Environment variable loaded: %s", key)
	}
	return value
}
