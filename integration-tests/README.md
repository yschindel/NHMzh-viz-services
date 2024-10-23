# Integration Tests

## Prerequisites

- Docker
- Docker Compose
- ts-node

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
docker compose up --build
```

### Open the MinIO console

In your browser, go to `localhost:9001` to see the MinIO console.

- Delete all items in the `ifc-files` bucket.
- Delete the `ifc-files` bucket.

### Add content to the Kafka topic

Open a terminal in the `integration-tests` directory.

```bash
ts-node test.ts
```

### Verify the content was added correctly

Verify that the bucket was created and that the content was added correctly:

1. Refresh the Object Browser in the MinIO console.
2. You should see a new bucket called `ifc-fragment-files`.
3. Open the `ifc-fragment-files` bucket. If no folders show up, try refreshing the bucket.
4. Open the folders and check that the files are present.

If you see the files, the integration tests passed.
