# Integration Tests

## Prerequisites

- Docker
- Docker Compose
- ts-node
- You have followed the instructions to set up the .env file in the [main README](../README.md)

### Install dependencies

Open a terminal in the `integration-tests` directory.

Install ts-node:

```bash
npm install -g ts-node
```

Install dependencies:

```bash
npm install
```

## Adding an IFC file to the test

1. Create a new directory called `assets` in the `integration-tests` directory.
2. Add an IFC file to the `assets` directory.
3. Rename the file to `test.ifc`.

## Running the tests

Open a new terminal window in the `integration-tests` directory:

```bash
cd integration-tests
```

### Start the services

Run the following command to start the services:

```bash
docker compose up --build -d
```

### Open the MinIO console

In your browser, go to `localhost:9001` to see the MinIO console.

- Login using the credentials you specified in the .env ( [main README](../README.md) )
- Delete all items in the `ifc-files` bucket.
- Delete the `ifc-files` bucket.
- Delete all items in the `ifc-fragment-files` bucket.
- Delete the `ifc-fragment-files` bucket.

### Add content to the Kafka topic

Open a terminal in the `integration-tests` directory and run:

```bash
ts-node test.ts
```

### Verify the content was added correctly

Verify that the bucket was created and that the content was added correctly:

1. Refresh the Object Browser in the MinIO console.
2. You should see new buckets called `ifc-files` and `ifc-fragment-files`.
3. Open the buckets. If no folders show up, try refreshing the bucket.
4. Open the folders and check that the files are present.

### Verify the viz_pbi-server is working

Check the output of the test.ts in the terminal. It should say: "Pass: Fragments file is a Buffer".

## MongoDB Compass

TODO: add description