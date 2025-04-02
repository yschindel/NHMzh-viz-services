package azure

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"viz_pbi-server/models"

	mssql "github.com/microsoft/go-mssqldb"
)

// SqlWriter handles writing messages to the database
type SqlWriter struct {
	db *sql.DB
}

// NewSqlWriter creates a new message writer
func NewSqlWriter(db *sql.DB) *SqlWriter {
	return &SqlWriter{
		db: db,
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
			log.Printf("Deadlock detected during %s (attempt %d of %d), retrying...",
				operation, attempt+1, maxRetries)
			time.Sleep(time.Millisecond * 100 * time.Duration(attempt+1))
			continue
		}

		return err
	}
	return fmt.Errorf("failed to %s after %d attempts", operation, maxRetries)
}

func (w *SqlWriter) writeElementsWithRetry(items []models.EavElementDataItem) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for data_eav table
	eavStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO data_eav_elements (project, filename, fileid, timestamp, id, param_name, param_value_string, param_value_number, param_value_boolean, param_value_date, param_type)
			VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11);`)
	if err != nil {
		return fmt.Errorf("error preparing data_eav statement: %v", err)
	}
	defer eavStmt.Close()

	for _, item := range items {
		// Write to data_eav table
		_, err = eavStmt.ExecContext(ctx,
			item.Project,           // @p1
			item.Filename,          // @p2
			item.FileID,            // @p3
			item.Timestamp,         // @p4
			item.Id,                // @p5
			item.ParamName,         // @p6
			item.ParamValueString,  // @p7
			item.ParamValueNumber,  // @p8
			item.ParamValueBoolean, // @p9
			item.ParamValueDate,    // @p10
			item.ParamType,         // @p11
		)
		if err != nil {
			return fmt.Errorf("error inserting data_eav_elements record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (w *SqlWriter) writeMaterialsWithRetry(items []models.EavMaterialDataItem) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for data_eav table
	eavStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO data_eav_materials (project, filename, fileid, timestamp, id, sequence, param_name, param_value_string, param_value_number, param_value_boolean, param_value_date, param_type)
			VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12);`)
	if err != nil {
		return fmt.Errorf("error preparing data_eav statement: %v", err)
	}
	defer eavStmt.Close()

	for _, item := range items {
		// Write to data_eav table
		_, err = eavStmt.ExecContext(ctx,
			item.Project,           // @p1
			item.Filename,          // @p2
			item.FileID,            // @p3
			item.Timestamp,         // @p4
			item.Id,                // @p5
			item.Sequence,          // @p6
			item.ParamName,         // @p7
			item.ParamValueString,  // @p8
			item.ParamValueNumber,  // @p9
			item.ParamValueBoolean, // @p10
			item.ParamValueDate,    // @p11
			item.ParamType,         // @p12
		)
		if err != nil {
			return fmt.Errorf("error inserting data_eav_materials record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (w *SqlWriter) WriteBlobData(item models.BlobData) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for data_updates table
	updateStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO data_updates (project, filename, timestamp, model_blob_url)
		VALUES (@p1, @p2, @p3, @p4);`)
	if err != nil {
		return fmt.Errorf("error preparing data_updates statement: %v", err)
	}
	defer updateStmt.Close()

	_, err = updateStmt.ExecContext(ctx,
		item.Project,   // @p1
		item.Filename,  // @p2
		item.Timestamp, // @p3
		item.BlobID,    // @p4
	)
	if err != nil {
		return fmt.Errorf("error inserting data_updates record: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
