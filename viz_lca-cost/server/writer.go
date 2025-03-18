package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

// MessageWriter handles writing messages to the database
type MessageWriter struct {
	db *sql.DB
}

// NewMessageWriter creates a new message writer
func NewMessageWriter(db *sql.DB) *MessageWriter {
	return &MessageWriter{
		db: db,
	}
}

// WriteLcaMessage writes an LCA message to the database
func (w *MessageWriter) WriteLcaMessage(message LcaMessage) error {
	return w.retryOnDeadlock("write lca message", func() error {
		return w.writeLcaMessageWithRetry(message)
	})
}

// WriteCostMessage writes a Cost message to the database
func (w *MessageWriter) WriteCostMessage(message CostMessage) error {
	return w.retryOnDeadlock("write cost message", func() error {
		return w.writeCostMessageWithRetry(message)
	})
}

// retryOnDeadlock retries the given function if a deadlock is detected
func (w *MessageWriter) retryOnDeadlock(operation string, fn func() error) error {
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

func (w *MessageWriter) writeCostMessageWithRetry(message CostMessage) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for cost_data table
	costStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO cost_data (project, filename, fileid, timestamp, id, category, level, is_structural, fire_rating, ebkph, ebkph_1, ebkph_2, ebkph_3, cost, cost_unit)
			VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15);`)
	if err != nil {
		return fmt.Errorf("error preparing cost_data statement: %v", err)
	}
	defer costStmt.Close()

	// Use FileID that's already been prepared by the processor
	fileID := message.FileID

	for _, item := range message.Data {
		// Write to cost_data table
		_, err = costStmt.ExecContext(ctx,
			message.Project,   // @p1
			message.Filename,  // @p2
			fileID,            // @p3
			message.Timestamp, // @p4
			item.Id,           // @p5
			item.Category,     // @p6
			item.Level,        // @p7
			item.IsStructural, // @p8
			item.FireRating,   // @p9
			item.Ebkph,        // @p10
			item.Ebkph1,       // @p11
			item.Ebkph2,       // @p12
			item.Ebkph3,       // @p13
			item.Cost,         // @p14
			item.CostUnit,     // @p15
		)
		if err != nil {
			return fmt.Errorf("error inserting cost_data record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (w *MessageWriter) writeLcaMessageWithRetry(message LcaMessage) error {
	// pretty print the message
	msgJson, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	log.Printf("LCA Message:\n%s", string(msgJson))

	ctx := context.Background()

	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the INSERT statement for lca_data table based on create_tables.sql
	lcaStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO lca_data (project, filename, fileid, timestamp, id, mat_kbob,
								  gwp_absolute, gwp_relative, penr_absolute, penr_relative,
								  ubp_absolute, ubp_relative)
			VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12);`)
	if err != nil {
		return fmt.Errorf("error preparing lca_data statement: %v", err)
	}
	defer lcaStmt.Close()

	// Use FileID that's already been prepared by the processor
	fileID := message.FileID

	for _, item := range message.Data {
		// Write to lca_data table
		_, err = lcaStmt.ExecContext(ctx,
			message.Project,   // @p1
			message.Filename,  // @p2
			fileID,            // @p3
			message.Timestamp, // @p4
			item.Id,           // @p5
			item.MaterialKbob, // @p6
			item.GwpAbsolute,  // @p7
			item.GwpRelative,  // @p8
			item.PenrAbsolute, // @p9
			item.PenrRelative, // @p10
			item.UbpAbsolute,  // @p11
			item.UbpRelative,  // @p12
		)
		if err != nil {
			return fmt.Errorf("error inserting lca_data record: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
