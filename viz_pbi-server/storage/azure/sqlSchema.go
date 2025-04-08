package azure

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"viz_pbi-server/logger"
)

// ReadSQLFile reads the SQL file and returns its content as a string
func ReadSQLFile(filePath string) (string, error) {
	log := logger.New().WithFields(logger.Fields{
		"file_path": filePath,
		"operation": "read_sql_file",
	})

	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(logger.Fields{
			"error": err,
		}).Error("Error opening SQL file")
		return "", fmt.Errorf("error opening SQL file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.WithFields(logger.Fields{
			"error": err,
		}).Error("Error reading SQL file")
		return "", fmt.Errorf("error reading SQL file: %v", err)
	}

	return string(content), nil
}

// InitializeDatabase initializes the database schema using external SQL files
func InitializeDatabase(db *sql.DB) error {
	baseLogger := logger.New().WithFields(logger.Fields{
		"operation": "initialize_database",
	})
	baseLogger.Info("Initializing database schema")

	ctx := context.Background()

	// Begin a transaction for all schema operations
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		baseLogger.WithFields(logger.Fields{
			"error": err,
		}).Error("Error starting schema transaction")
		return fmt.Errorf("error starting schema transaction: %v", err)
	}
	defer tx.Rollback()

	// Read and execute SQL files
	sqlFiles := []string{"/app/storage/azure/sql/create_tables.sql", "/app/storage/azure/sql/add_columns.sql"}
	for _, filePath := range sqlFiles {
		fileLogger := baseLogger.WithFields(logger.Fields{
			"file_path": filePath,
		})
		fileLogger.Info("Processing SQL file")

		sqlContent, err := ReadSQLFile(filePath)
		if err != nil {
			// Try the fallback path if the Docker path fails
			fileLogger.WithFields(logger.Fields{
				"error": err,
			}).Warn("Failed to read from Docker path, trying relative path")

			relativePath := fmt.Sprintf("./storage/azure/sql/%s", filePath[len("/app/storage/azure/sql/"):])
			sqlContent, err = ReadSQLFile(relativePath)
			if err != nil {
				fileLogger.WithFields(logger.Fields{
					"error":         err,
					"relative_path": relativePath,
				}).Error("Error reading SQL file (both Docker and relative paths)")
				return fmt.Errorf("error reading SQL file (both Docker and relative paths): %v", err)
			}
		}

		_, err = tx.ExecContext(ctx, sqlContent)
		if err != nil {
			fileLogger.WithFields(logger.Fields{
				"error": err,
			}).Error("Error executing SQL from file")
			return fmt.Errorf("error executing SQL from file %s: %v", filePath, err)
		}

		fileLogger.Info("Successfully executed SQL file")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		baseLogger.WithFields(logger.Fields{
			"error": err,
		}).Error("Error committing schema changes")
		return fmt.Errorf("error committing schema changes: %v", err)
	}

	baseLogger.Info("Database schema initialization completed successfully")
	return nil
}
