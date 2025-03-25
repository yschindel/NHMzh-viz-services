# Integration Tests

## Prerequisites

- Docker
- Docker Compose
- Go 1.20 or higher
- You have followed the instructions to set up the `.env` and `.env.dev` files in the [main README](../README.md)

### Install Dependencies

Open a terminal in the `integration-tests/lca_cost` directory.

Install dependencies:

```bash
go mod tidy
```

## Add an IFC File to the Test

1. Create a new directory called `assets` in the `integration-tests` directory.
2. Add an IFC file to the `assets` directory.
3. Rename the file to `test.ifc`.

## Test Setup

### Start the Azure SQL Server

The integration tests use the SQL Server container defined in the docker-compose.yml file. There's no need to use the Azure SQL Server Emulator in VSCode for running the integration tests.

### Start the Services

Open a new terminal window in the `integration-tests` directory:

```bash
cd integration-tests
docker compose up --build -d
```

This will start all the necessary services:

- MinIO (object storage)
- Kafka and Zookeeper (message broker)
- SQL Server (database)
- SQL Server initialization (creates the required database)
- viz_ifc (IFC processing service)
- viz_lca-cost (cost and LCA data service)
- viz_pbi-server (Power BI integration server)

You should now see all the containers running for the integration tests services.

## Running the IFC Integration Tests

### Open the MinIO Console

In your browser, go to `localhost:9001` to see the MinIO console.

- Login using the credentials you specified in the `.env` file (default: ROOTUSER/CHANGEME123)
- Before running tests, you may want to delete any existing data:
  - Delete all items in the `ifc-files` bucket and the bucket itself
  - Delete all items in the `ifc-fragment-files` bucket and the bucket itself
  - Delete all items in the `lca-cost-data` bucket and the bucket itself

### Add Content to the IFC Kafka Topic

Open a terminal in the `integration-tests/ifc` directory and run:

```bash
go run main.go
```

Watch the output of the terminal. It should tell you that the test was passed. This takes up to a minute.

**Verify the content was added correctly:**

Verify that the bucket was created and that the content was added correctly:

1. Refresh the Object Browser in the MinIO console.
2. You should see new buckets called `ifc-files` and `ifc-fragment-files`.
3. Open the buckets. If no folders show up, try refreshing the bucket.
4. Open the folders and check that the files are present.

### Verify the viz_pbi-server is Working

Check the output of the test.ts in the terminal. It should say: "Pass: Fragments file is a Buffer".

## Running the LCA and Cost Integration Tests

### Add Content to the LCA and Cost Kafka Topic

Open a terminal in the `integration-tests/lca_cost` directory and run:

```bash
go run main.go
```

**Verify the content was added correctly:**

Connect to the SQL Server and check that the data has been added. You can use:

1. The SQL Server tab in VSCode:

   - Connect to `localhost,1433` with username `sa` and password `P@ssw0rd`
   - Navigate to the `nhmzh-viz-services-localdev` database
   - Check that the tables have been created and populated

2. Or use a SQL client like SQL Server Management Studio or Azure Data Studio:
   - Connect to `localhost,1433` with username `sa` and password `P@ssw0rd`
   - Run a query to verify the data: `SELECT TOP 10 * FROM [nhmzh-viz-services-localdev].[dbo].[LcaData]`

## Troubleshooting

### SQL Server Connection Issues

If you have issues connecting to SQL Server:

- Ensure the SQL Server container is running (`docker ps`)
- Check if the database was created successfully in the logs (`docker logs sqlserver-init`)
- Verify you're using the correct connection string parameters

### Kafka Issues

If you have issues with Kafka:

- Check if the Kafka container is running (`docker ps`)
- Verify the topics are created (`docker exec -it viz-test-kafka-container kafka-topics.sh --list --bootstrap-server kafka:9093`)

### MinIO Issues

If you have issues with MinIO:

- Check if the MinIO container is running (`docker ps`)
- Verify you can access the MinIO console at http://localhost:9001
- Check if buckets are created and accessible
