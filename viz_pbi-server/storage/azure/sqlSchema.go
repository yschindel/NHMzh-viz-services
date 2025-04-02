package azure

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
)

// ReadSQLFile reads the SQL file and returns its content as a string
func ReadSQLFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening SQL file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading SQL file: %v", err)
	}

	return string(content), nil
}

// InitializeDatabase initializes the database schema using external SQL files
func InitializeDatabase(db *sql.DB) error {
	log.Println("Initializing database schema...")

	ctx := context.Background()

	// Begin a transaction for all schema operations
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting schema transaction: %v", err)
	}
	defer tx.Rollback()

	// Read and execute SQL files
	sqlFiles := []string{"/app/storage/azure/sql/create_tables.sql", "/app/storage/azure/sql/add_columns.sql"}
	for _, filePath := range sqlFiles {
		log.Printf("Attempting to read SQL file: %s", filePath)
		sqlContent, err := ReadSQLFile(filePath)
		if err != nil {
			// Try the fallback path if the Docker path fails
			log.Printf("Failed to read from Docker path, trying relative path: %v", err)
			relativePath := fmt.Sprintf("./storage/azure/sql/%s", filePath[len("/app/storage/azure/sql/"):])
			sqlContent, err = ReadSQLFile(relativePath)
			if err != nil {
				return fmt.Errorf("error reading SQL file (both Docker and relative paths): %v", err)
			}
		}

		_, err = tx.ExecContext(ctx, sqlContent)
		if err != nil {
			return fmt.Errorf("error executing SQL from file %s: %v", filePath, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing schema changes: %v", err)
	}

	log.Println("Database schema initialization completed successfully")
	return nil
}
