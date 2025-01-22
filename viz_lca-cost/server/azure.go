package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
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

func retryOnDeadlock(operation string, fn func() error) error {
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		// Check for SQL Server deadlock error (1205)
		if sqlErr, ok := err.(mssql.Error); ok && sqlErr.Number == 1205 {
			log.Printf("Deadlock detected during %s (attempt %d of %d), retrying...",
				operation, attempt+1, maxRetries)
			time.Sleep(time.Millisecond * 100 * time.Duration(attempt+1))
			continue
		}

		return err
	}
	return fmt.Errorf("failed to %s after %d attempts", operation, maxRetries)
}

func WriteLcaMessage(db *sql.DB, message LcaMessage) error {
	return retryOnDeadlock("write lca message", func() error {
		return writeLcaMessageWithRetry(db, message)
	})
}

func WriteCostMessage(db *sql.DB, message CostMessage) error {
	return retryOnDeadlock("write cost message", func() error {
		return writeCostMessageWithRetry(db, message)
	})
}

func writeCostMessageWithRetry(db *sql.DB, message CostMessage) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
			MERGE project_data AS target
			USING (VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)) 
					AS source (project, filename, timestamp, id, category, cost, cost_unit)
			ON target.id = source.id 
			   AND target.project = source.project 
			   AND target.filename = source.filename
			WHEN MATCHED THEN
					UPDATE SET 
							target.timestamp = source.timestamp,
							target.cost = source.cost,
							target.cost_unit = source.cost_unit
			WHEN NOT MATCHED THEN
					INSERT (project, filename, timestamp, id, category, cost, cost_unit)
					VALUES (source.project, source.filename, source.timestamp, source.id, 
								 source.category, source.cost, source.cost_unit);`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, item := range message.Data {
		_, err = stmt.ExecContext(ctx,
			message.Project,
			message.Filename,
			message.Timestamp,
			item.Id,
			item.Category,
			item.Cost,
			item.CostUnit,
		)
		if err != nil {
			return fmt.Errorf("error merging cost record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func writeLcaMessageWithRetry(db *sql.DB, message LcaMessage) error {
	// pretty print the message
	msgJson, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	log.Printf("LCA Message:\n%s", string(msgJson))

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
			MERGE project_data AS target
			USING (VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12)) 
					AS source (project, filename, timestamp, id, category, material_kbob, 
										gwp_absolute, gwp_relative, penr_absolute, penr_relative, 
										ubp_absolute, ubp_relative)
			ON target.id = source.id 
			   AND target.project = source.project 
			   AND target.filename = source.filename
			WHEN MATCHED THEN
					UPDATE SET 
							target.timestamp = source.timestamp,
							target.category = source.category,
							target.material_kbob = source.material_kbob,
							target.gwp_absolute = source.gwp_absolute,
							target.gwp_relative = source.gwp_relative,
							target.penr_absolute = source.penr_absolute,
							target.penr_relative = source.penr_relative,
							target.ubp_absolute = source.ubp_absolute,
							target.ubp_relative = source.ubp_relative
			WHEN NOT MATCHED THEN
					INSERT (project, filename, timestamp, id, category, material_kbob,
								 gwp_absolute, gwp_relative, penr_absolute, penr_relative,
								 ubp_absolute, ubp_relative)
					VALUES (source.project, source.filename, source.timestamp, source.id,
								 source.category, source.material_kbob, source.gwp_absolute,
								 source.gwp_relative, source.penr_absolute, source.penr_relative,
								 source.ubp_absolute, source.ubp_relative);`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, item := range message.Data {
		_, err = stmt.ExecContext(ctx,
			message.Project,
			message.Filename,
			message.Timestamp,
			item.Id,
			item.Category,
			item.MaterialKbob,
			item.GwpAbsolute,
			item.GwpRelative,
			item.PenrAbsolute,
			item.PenrRelative,
			item.UbpAbsolute,
			item.UbpRelative,
		)
		if err != nil {
			return fmt.Errorf("error merging lca record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
