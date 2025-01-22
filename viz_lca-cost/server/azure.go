package server

import (
	"context"
	"database/sql"
	"fmt"

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
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=false;trustservercertificate=true;",
		config.Server, config.User, config.Password, config.Port, config.Database)

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
func WriteCostMessage(db *sql.DB, message CostMessage) error {
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
			WHEN MATCHED THEN
					UPDATE SET 
							target.project = source.project,
							target.filename = source.filename,
							target.timestamp = source.timestamp,
							target.category = source.category,
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

func WriteLcaMessage(db *sql.DB, message LcaMessage) error {
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
			WHEN MATCHED THEN
					UPDATE SET 
							target.project = source.project,
							target.filename = source.filename,
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
