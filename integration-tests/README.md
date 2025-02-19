# Integration Tests

## Prerequisites

- Docker
- Docker Compose
- Go 1.23.2 (or higher)
- You have followed the instructions to set up the .env file in the [main README](../README.md)

### Install dependencies

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

### Start the Azure SQL Server Emulator

Use the 'mssql' extension in VSCode to start the Azure SQL Server Emulator.
Follow the instructions in this video closely: https://www.youtube.com/watch?v=3XgepwpBJP8

Basic steps:
Go to the Database Projects section of the 'mssql' extension in VSCode.
Right click on the project and select 'Build'. >> follow the video tutorial
Right click on the project and select 'Publish'. >> follow the video tutorial

Switch to the 'SQL Server' tab in VSCode.
Connect to the Azure SQL Server Emulator. Skip all optionals inputs, use 'sa' as the username and provide a password.

Open the project_data.sql file.
In the bottom right corner of VSCode where it says 'Disconnect', click on it and connect to the emulator. make sure to select the 'sa' user and later the right 'database' or whatever it is called. Choose (nhmzh-viz-services-localdev) as the database.
Run 'project_data.sql' to create the database table.

Open a new terminal window in the `integration-tests` directory:

```bash
cd integration-tests
```

### Start the services

Run the following command to start the services:

```bash
docker compose up --build -d
```

You should now have one container running for the Azure SQL Server Emulator and a contaier group (or however that is called) running for all the integration tests services.

## Running the IFC Integration Tests

### Open the MinIO Console

In your browser, go to `localhost:9001` to see the MinIO console.

- Login using the credentials you specified in the .env ( [main README](../README.md) )
- Delete all items in the `ifc-files` bucket.
- Delete the `ifc-files` bucket.
- Delete all items in the `ifc-fragment-files` bucket.
- Delete the `ifc-fragment-files` bucket.
- Delete all items in the `lca-cost-data` bucket.
- Delete the `lca-cost-data` bucket.

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

### Add content to the LCA and Cost Kafka topic

Open a terminal in the `integration-tests/lca_cost` directory and run:

```bash
go run main.go
```

**Verify the content was added correctly:**

In the (VS Code) environment where you set up the Azure SQL Server Emulator, go to the SQL Server Explorer.
