# Integration Tests

## Prerequisites

- Docker
- Docker Compose
- Go 1.23.2 (or higher)
- You have followed the instructions to set up the .env file in the [main README](../README.md)
- MongoDB Compass (Desktop viewer for MongoDB, optional)

### Install dependencies

Open a terminal in the `integration-tests/lca` directory.

Install dependencies:

```bash
go mod tidy
```

Do the same in the `integration-tests/ifc` directory.

## Add an IFC File to the Test

1. Create a new directory called `assets` in the `integration-tests` directory.
2. Add an IFC file to the `assets` directory.
3. Rename the file to `test.ifc`.

## Running the Tests

Open a new terminal window in the `integration-tests` directory:

```bash
cd integration-tests
```

### Start the services

Run the following command to start the services:

```bash
docker compose up --build -d
```

### Open the MinIO Console

In your browser, go to `localhost:9001` to see the MinIO console.

- Login using the credentials you specified in the .env ( [main README](../README.md) )
- Delete all items in the `ifc-files` bucket.
- Delete the `ifc-files` bucket.
- Delete all items in the `ifc-fragment-files` bucket.
- Delete the `ifc-fragment-files` bucket.

### Add Content to the IFC Kafka Topic

Open a terminal in the `integration-tests/ifc` directory and run:

```bash
go run main.go
```

**Verify the content was added correctly:**

Verify that the bucket was created and that the content was added correctly:

1. Refresh the Object Browser in the MinIO console.
2. You should see new buckets called `ifc-files` and `ifc-fragment-files`.
3. Open the buckets. If no folders show up, try refreshing the bucket.
4. Open the folders and check that the files are present.

### Add content to the LCA Kafka topic

Open a terminal in the `integration-tests/lca_cost` directory and run:

```bash
go run main.go
```

**Verify the content was added correctly:**

1. Open MongoDB Compass
2. Check the 'testdb' Database
3. Check the 'testcollection' in testdb
4. You should see 3 documents

### Verify the viz_pbi-server is Working

Check the output of the test.ts in the terminal. It should say: "Pass: Fragments file is a Buffer".

