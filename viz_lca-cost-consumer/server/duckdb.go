package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type DuckDBManager struct {
	db          *sql.DB
	minioClient *minio.Client
	credentials CustomMinioCredentials
	bucket      string
	mutex       sync.Mutex
}

func NewDuckDBManager(dbPath string, customCredentials CustomMinioCredentials, bucket string) (*DuckDBManager, error) {
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open duckdb: %w", err)
	}

	// Initialize Minio client
	minioClient, err := minio.New(customCredentials.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(customCredentials.AccessKeyID, customCredentials.SecretAccessKey, ""),
		Secure: false, // Set to true if using HTTPS
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	manager := &DuckDBManager{
		db:          db,
		minioClient: minioClient,
		credentials: customCredentials,
		bucket:      bucket,
	}

	// Initialize S3 connection for DuckDB
	if err := manager.initS3Connection(); err != nil {
		db.Close()
		return nil, err
	}

	return manager, nil
}

// Add this method to the DuckDBManager struct
func (m *DuckDBManager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *DuckDBManager) initS3Connection() error {
	_, err := m.db.Exec(`
		INSTALL httpfs;
		LOAD httpfs;
		INSTALL s3;
		LOAD s3;
	`)
	if err != nil {
		return fmt.Errorf("failed to load extensions: %w", err)
	}

	_, err = m.db.Exec(fmt.Sprintf(`
		SET s3_access_key_id='%s';
		SET s3_secret_access_key='%s';
		SET s3_endpoint='%s';
		SET s3_use_ssl=false;
		SET s3_url_style='path';
	`, m.credentials.AccessKeyID, m.credentials.SecretAccessKey, m.credentials.Endpoint))

	return err
}

func (m *DuckDBManager) ensureBucket() error {
	exists, err := m.minioClient.BucketExists(context.Background(), m.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = m.minioClient.MakeBucket(context.Background(), m.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket: %s", m.bucket)
	}

	return nil
}

// EnsureParquetFile checks if the parquet file exists and creates it if it doesn't
func (m *DuckDBManager) EnsureParquetFile(project string, filename string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// First ensure the bucket exists using Minio client
	if err := m.ensureBucket(); err != nil {
		return err
	}

	objectPath := fmt.Sprintf("%s/%s.parquet", project, filename)

	// Try to read the file to check if it exists
	checkQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM read_parquet('s3://%s/%s');
	`, m.bucket, objectPath)

	_, err := m.db.Exec(checkQuery)
	if err != nil {
		// File doesn't exist, create it with initial schema
		createQuery := fmt.Sprintf(`
			COPY (
				SELECT 
					id,
					NULL::VARCHAR as category,
					NULL::FLOAT as co2e,
					NULL::FLOAT as greyEnergy,
					NULL::FLOAT as UBP,
					NULL::FLOAT as cost
				FROM (SELECT '' as id) WHERE false
			) TO 's3://%s/%s'
			(FORMAT 'parquet', COMPRESSION 'zstd')
		`, m.bucket, objectPath)

		_, err = m.db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create initial parquet file: %w", err)
		}
		log.Printf("Created parquet file: s3://%s/%s", m.bucket, objectPath)
	}

	return nil
}

// UpdateParquetFromEnvironmentalData updates the parquet file with environmental impact data
func (m *DuckDBManager) UpdateParquetFromEnvironmentalData(project, filename string, data []LcaDataItem) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	objectPath := fmt.Sprintf("%s/%s.parquet", project, filename)

	// write json to disk
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	os.WriteFile("new_data.json", jsonData, 0644)

	_, err = m.db.Exec(`
		DROP TABLE IF EXISTS new_lca_data; 
		CREATE TABLE new_lca_data (id VARCHAR, category VARCHAR, co2e FLOAT, greyEnergy FLOAT, UBP FLOAT); 
		INSERT INTO new_lca_data SELECT * FROM read_json('new_data.json', columns = {id: 'VARCHAR', category: 'VARCHAR', co2e: 'FLOAT', greyEnergy: 'FLOAT', UBP: 'FLOAT'})
	`)
	if err != nil {
		return fmt.Errorf("failed to insert data into temp table: %w", err)
	}

	// Merge with existing data and write back to parquet
	mergeQuery := fmt.Sprintf(`
		COPY (
			WITH existing AS (
				SELECT * FROM read_parquet('s3://%s/%s')
			),
			merged AS (
				SELECT 
					COALESCE(n.id, e.id) as id,
					COALESCE(n.category, e.category) as category,
					COALESCE(n.co2e, e.co2e) as co2e,
					COALESCE(n.greyEnergy, e.greyEnergy) as greyEnergy,
					COALESCE(n.UBP, e.UBP) as UBP,
					e.cost
				FROM existing e
				FULL OUTER JOIN new_lca_data n ON e.id = n.id
			)
			SELECT * FROM merged
		)
		TO 's3://%s/%s'
		(FORMAT 'parquet', COMPRESSION 'zstd')
	`, m.bucket, objectPath, m.bucket, objectPath)

	_, err = m.db.Exec(mergeQuery)
	if err != nil {
		return fmt.Errorf("failed to merge data: %w", err)
	}

	return err
}

// UpdateParquetFromCostData updates the parquet file with cost data
func (m *DuckDBManager) UpdateParquetFromCostData(project, filename string, data []CostDataItem) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	objectPath := fmt.Sprintf("%s/%s.parquet", project, filename)

	// write json to disk
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	os.WriteFile("new_cost_data.json", jsonData, 0644)

	_, err = m.db.Exec(`
		DROP TABLE IF EXISTS new_cost_data; 
		CREATE TABLE new_cost_data (id VARCHAR, category VARCHAR, cost FLOAT); 
		INSERT INTO new_cost_data SELECT * FROM read_json('new_cost_data.json', columns = {id: 'VARCHAR', category: 'VARCHAR', cost: 'FLOAT'})
	`)
	if err != nil {
		return fmt.Errorf("failed to insert data into temp table: %w", err)
	}

	// Merge with existing data and write back to parquet
	mergeQuery := fmt.Sprintf(`
		COPY (
			WITH existing AS (
				SELECT * FROM read_parquet('s3://%s/%s')
			),
			merged AS (
				SELECT 
					COALESCE(n.id, e.id) as id,
					e.category,
					e.co2e,
					e.greyEnergy,
					e.UBP,
					COALESCE(n.cost, e.cost) as cost
				FROM existing e
				FULL OUTER JOIN new_cost_data n ON e.id = n.id
			)
			SELECT * FROM merged
		)
		TO 's3://%s/%s'
		(FORMAT 'parquet', COMPRESSION 'zstd')
	`, m.bucket, objectPath, m.bucket, objectPath)

	_, err = m.db.Exec(mergeQuery)

	return err
}
