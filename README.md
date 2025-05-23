# NHMzh Visualization Services

This repository contains the code for the Visualization Services of the NHMzh project.

## The Services in this Repository

### Minio

- Used for storing files.
- Buckets:
  - `ifc-files` (the raw IFC files)
  - `ifc-fragment-files` (compressed counterparts to the IFC files, converted to 'fragments')

### viz_ifc

Listens to a Kafka topic with links to IFC files:

- Loads the IFC file from Minio
- Converts to fragments, compresses
- Saves compressed fragments back to Minio

Uses the @ThatOpen Companies library.

### viz_lca-cost

This consumer listens to the cost and LCA Kafka topics and writes data to Azure SQL DB:

- Cost and LCA data is captured in Azure SQL DB
- This enables reporting of the data over time

### viz_pbi-server

Serves the Power BI integration:

- Provides an API for accessing fragment files stored in Minio

## Additional Services Needed:

- Azure SQL Server SQL Database for storing all cost and LCA data over time.
- These services can be connected to PowerBI via direct query so it always shows the latest data without having to update the dashboard in PowerBI desktop and republish.

## Prerequisites

- Docker
- Docker Compose
- Node.js
- Go
- Azure SQL Server Emulator (VSCode extension 'mssql')

## Environment Variables

Create a `.env` file in the root directory of the repository with the following variables:

```
ENVIRONMENT=production

MINIO_ACCESS_KEY=ROOTUSER
MINIO_SECRET_KEY=CHANGEME123

MINIO_ENDPOINT=minio
MINIO_PORT=9000

MINIO_USE_SSL=false

MINIO_IFC_BUCKET=ifc-files
VIZ_FRAGMENTS_BUCKET=ifc-fragment-files
MINIO_LCA_COST_DATA_BUCKET=lca-cost-data

VIZ_KAFKA_IFC_CONSUMER_ID=viz-ifc-consumer
VIZ_KAFKA_IFC_GROUP_ID=viz-ifc
VIZ_KAFKA_DATA_GROUP_ID=viz-data

PBI_SERVER_PORT=3000

KAFKA_BROKER=kafka:9093
KAFKA_IFC_TOPIC=ifc-files
KAFKA_LCA_TOPIC=lca-data
KAFKA_COST_TOPIC=cost-data

# Azure DB credentials
AZURE_DB_SERVER=your-server-name.database.windows.net
AZURE_DB_PORT=1433
AZURE_DB_USER=your-user-name
AZURE_DB_PASSWORD=your-password
AZURE_DB_DATABASE=your-database-name
```

If you intend to do local integration testing:
Create a `.env.dev` file next to the `.env` file and override the Azure DB credentials:

```
# Azure DB credentials
AZURE_DB_SERVER=host.docker.internal
AZURE_DB_PORT=1234
AZURE_DB_USER=your-user-name
AZURE_DB_PASSWORD=your-password
AZURE_DB_DATABASE=your-database-name

ENVIRONMENT=development
```

Make sure to replace `ROOTUSER` and `CHANGEME123` with your own credentials.

## Running the Services

To run the services, navigate to the root directory of the repository and run the following command:

```bash
docker-compose up --build -d
```

## Stopping the Services

To stop the services, navigate to the root directory of the repository and run the following command:

```bash
docker-compose down
```

## Access the MinIO Console

To access the MinIO Console, navigate to http://localhost:9001. You will need to use the credentials `ROOTUSER` and `CHANGEME123` to log in (or whatever you have set in the `.env` file).

## Tests

### Integration Tests

Follow the instructions in the [integration-tests/README.md](integration-tests/README.md) file.
