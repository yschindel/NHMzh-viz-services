package azure

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/microsoft/go-mssqldb"
)

type DBConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

func ConnectDB(config DBConfig) (*sql.DB, error) {
	environment := os.Getenv("ENVIRONMENT")

	// Build connection string with environment-specific settings
	var connString string
	if environment == "development" {
		// Development connection string with SSL disabled
		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=false;trustservercertificate=true;",
			config.Server, config.User, config.Password, config.Port, config.Database)
	} else {
		// Production connection string with SSL enabled
		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
			config.Server, config.User, config.Password, config.Port, config.Database)
	}

	log.Printf("Connecting to database in %s mode", environment)

	// Create connection pool
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err)
	}

	// Test connection
	err = db.PingContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error testing connection: %v", err)
	}

	return db, nil
}
