# Integration Tests

Testing the backend services.

## Overview

The integration tests are designed to simulate a full run from IFC upload all the way to Azure SQL and Blob Storage.
They are executed by running a .go file that uploads a file to MinIO and then sends a message with the download link to that file through a Kafka topic. All services to their processing and after some time the final data should be availabe in the Azure SQL DB and Blob Storage. Depeding on the configuration of environment variables, this can be local emulators or a real Azure instance.

## Prerequisites

- Docker
- Docker Compose
- Go 1.20 or higher
- You have followed the instructions to set up the `.env` and `.env.dev` files in the [main README](../README.md)

## Test Setup

Before we can start the services, we need to set all environment varialbes in the docker-compose.yml file.

### Add an IFC File to the Test

1. Create a new directory called `assets` in the `integration-tests` directory.
2. Add an IFC file to the `assets` directory.
3. Rename the file to `test.ifc`.

### Install Dependencies

Open a terminal in the `integration-tests/test` directory.

Install dependencies:

```bash
go mod tidy
```

### Start the Services

Open a new terminal window in the `integration-tests` directory:

```bash
cd integration-tests
docker compose up --build -d
```

This will start all the necessary services:

- MinIO
  - Used to upload a test IFC file to
- Kafka and Zookeeper
  - Messageging system
  - Topics used in test (not the actual topic names):
    - New ifc file available
    - Lca data
    - Cost data
- SQL Server
  - Azure emulator
  - For data Storage
- SQL Server initialization (creates the required database)
- Azurite
  - Azure emulator
  - Blob storage
- mock_calc_lca-cost
  - Mocks calculations performed by the actual services
  - Lca data
  - Cost data
- viz_ifc (IFC processing service)
  - Processing IFC files
  - Converts IFC file to fragments for WebGL viewing
  - Extracts general element data from IFC files
- viz_lca-cost (cost and LCA data service)
  - Reshaping data from 'row' or 'object' style to EAV
  - Forwards to client specific storage service
- viz_pbi-server
  - Client specific storage service
  - Backend for PowerBI 3D viewer fragments loading
  - Single point of communication with Azure (Except Direct Query from PowerBI)

You should now see all the containers running for the integration tests services.

## Running the IFC Integration Tests

### Open the MinIO Console (Optional)

In your browser, go to `localhost:9001` to see the MinIO console.

- Login using the credentials you specified in the `.env` file (default: ROOTUSER/CHANGEME123)
- Before running tests, you may want to delete any existing data (this is optional):
  - Delete all items in the `ifc-files` bucket and the bucket itself

### Trigger the Test

Open a terminal in the `integration-tests/test` directory and run:

```bash
go run main.go
```

Watch the output of the terminal. It should tell you that the test was passed. This takes up to a minute or two.

Here's what you need to watch out for to verify the test passed:

**Automatically checked:**

The test will automatically check if the expected fragments file can be retrived from the pbi-server. It will log success or failure on that in the terminal.

**Verify the ifc file was added correctly:**

Verify that the ifc bucket was created and that the file was added correctly:

1. Refresh the Object Browser in the MinIO console.
2. You should see a new bucket called `ifc-files`.
3. Open the buckets. If no folders show up, try refreshing the bucket.
4. Open the folders and check that the files are present.

**Verify the data was added to Azure:**

Depending on you environment varialbes, connect to the local emulators or the actual cloud instance.
This can be done with the 'mssql' extension for vs code.

You should see 3 tables:

- element data
  - EAV (element-attribute-value) style
  - Contains element-level information
  - Column 'param_name' should include what you have specified in the environment variable for the viz_ifc service (if it exists in the IFC model):
    - Category
    - Level
    - LoadBearing
    - ...
- material data
  - EAV style
  - Contains material-level information
    - Material 'layers' are grouped by 'sequence'
      - Every 'sequence' should have the same set of 'param_name' (6 or so)
    - Multiple materials per 'id' aka element
- updates
  - One row per new IFC upload
  - Metadata about the data update
