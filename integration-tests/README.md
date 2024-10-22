# Integration Tests

## Prerequisites

- Docker
- Docker Compose

## Running the tests

Make sure you are in the `integration-tests` directory:

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

Run the following command to add content to the Kafka topic:

```bash
bash ./testProduceConsume.bash
```

### Verify the content was added correctly

Verify that the bucket was created and that the content was added correctly.
