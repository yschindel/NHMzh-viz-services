package azure

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"viz_pbi-server/logger"
	"viz_pbi-server/models"

	mssql "github.com/microsoft/go-mssqldb"
)

// SqlWriter handles writing messages to the database
type SqlWriter struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewSqlWriter creates a new message writer
func NewSqlWriter(db *sql.DB) *SqlWriter {
	return &SqlWriter{
		db:     db,
		logger: logger.New(),
	}
}

// WriteMaterials writes an LCA message to the database
func (w *SqlWriter) WriteMaterials(dataItems []models.EavMaterialDataItem) error {
	return w.retryOnDeadlock("write lca message", func() error {
		return w.writeMaterialsWithRetry(dataItems)
	})
}

// WriteElements writes a Cost message to the database
func (w *SqlWriter) WriteElements(dataItems []models.EavElementDataItem) error {
	return w.retryOnDeadlock("write cost message", func() error {
		return w.writeElementsWithRetry(dataItems)
	})
}

// retryOnDeadlock retries the given function if a deadlock is detected
func (w *SqlWriter) retryOnDeadlock(operation string, fn func() error) error {
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		// Check for SQL Server deadlock error (1205)
		if sqlErr, ok := err.(mssql.Error); ok && sqlErr.Number == 1205 {
			w.logger.WithFields(logger.Fields{"error": err, "operation": operation, "attempt": attempt + 1, "maxRetries": maxRetries}).Warn("Deadlock detected, retrying...")
			time.Sleep(time.Millisecond * 100 * time.Duration(attempt+1))
			continue
		}

		return err
	}
	return fmt.Errorf("failed to %s after %d attempts", operation, maxRetries)
}

func (w *SqlWriter) writeElementsWithRetry(items []models.EavElementDataItem) error {
	ctx := context.Background()
	log := w.logger.WithFields(logger.Fields{
		"operation": "write_elements",
		"count":     len(items),
	})

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error starting transaction")
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for data_eav table
	eavStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO data_eav_elements (project, filename, timestamp, id, param_name, param_value_string, param_value_number, param_value_boolean, param_value_date, param_type)
			VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10);`)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error preparing data_eav statement")
		return fmt.Errorf("error preparing data_eav statement: %v", err)
	}
	defer eavStmt.Close()

	for _, item := range items {
		if item.ParamName == "" {
			log.WithFields(logger.Fields{"item": item}).Warn("Skipping item with empty param_name")
			continue
		}
		// Write to data_eav table
		_, err = eavStmt.ExecContext(ctx,
			item.Project,           // @p1
			item.Filename,          // @p2
			item.Timestamp,         // @p3
			item.Id,                // @p4
			item.ParamName,         // @p5
			item.ParamValueString,  // @p6
			item.ParamValueNumber,  // @p7
			item.ParamValueBoolean, // @p8
			item.ParamValueDate,    // @p9
			item.ParamType,         // @p10
		)
		if err != nil {
			log.WithFields(logger.Fields{"error": err}).Error("Error inserting data_eav_elements record")
			return fmt.Errorf("error inserting data_eav_elements record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error committing transaction")
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Info("Successfully wrote elements to database")
	return nil
}

func (w *SqlWriter) writeMaterialsWithRetry(items []models.EavMaterialDataItem) error {
	ctx := context.Background()
	log := w.logger.WithFields(logger.Fields{
		"operation": "write_materials",
		"count":     len(items),
	})

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error starting transaction")
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for data_eav table
	eavStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO data_eav_materials (project, filename, timestamp, id, sequence, param_name, param_value_string, param_value_number, param_value_boolean, param_value_date, param_type)
			VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11);`)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error preparing data_eav statement")
		return fmt.Errorf("error preparing data_eav statement: %v", err)
	}
	defer eavStmt.Close()

	for _, item := range items {
		// Write to data_eav table
		_, err = eavStmt.ExecContext(ctx,
			item.Project,           // @p1
			item.Filename,          // @p2
			item.Timestamp,         // @p3
			item.Id,                // @p4
			item.Sequence,          // @p5
			item.ParamName,         // @p6
			item.ParamValueString,  // @p7
			item.ParamValueNumber,  // @p8
			item.ParamValueBoolean, // @p9
			item.ParamValueDate,    // @p10
			item.ParamType,         // @p11
		)
		if err != nil {
			log.WithFields(logger.Fields{"error": err}).Error("Error inserting data_eav_materials record")
			return fmt.Errorf("error inserting data_eav_materials record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error committing transaction")
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Info("Successfully wrote materials to database")
	return nil
}

func (w *SqlWriter) WriteBlobData(item models.BlobData) error {
	ctx := context.Background()
	log := w.logger.WithFields(logger.Fields{
		"operation": "write_blob_data",
		"project":   item.Project,
		"filename":  item.Filename,
		"container": item.Container,
	})

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error starting transaction")
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for data_updates table
	updateStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO data_updates (project, filename, timestamp, model_blob_storage_url, model_blob_storage_container, model_blob_id)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6);`)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error preparing data_updates statement")
		return fmt.Errorf("error preparing data_updates statement: %v", err)
	}
	defer updateStmt.Close()

	_, err = updateStmt.ExecContext(ctx,
		item.Project,           // @p1
		item.Filename,          // @p2
		item.Timestamp,         // @p3
		item.StorageServiceURL, // @p4
		item.Container,         // @p5
		item.BlobID,            // @p6
	)
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error inserting data_updates record")
		return fmt.Errorf("error inserting data_updates record: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(logger.Fields{"error": err}).Error("Error committing transaction")
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Info("Successfully wrote blob data to database")
	return nil
}
