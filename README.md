# NHMzh Vizualization Services

This repository contains the code for the Vizualization Services of the NHMzh project.

## Prerequisites

- Docker
- Docker Compose
- Node.js (for running the unit tests locally)

## Environment Variables

Create a `.env` file in the root directory of the repository with the following variables:

```
MINIO_ROOT_USER=ROOTUSER
MINIO_ROOT_PASSWORD=CHANGEME123

MINIO_ENDPOINT=minio
MINIO_PORT=9000
MINIO_USE_SSL=false

MINIO_ACCESS_KEY=ROOTUSER
MINIO_SECRET_KEY=CHANGEME123

VIZ_KAFKA_IFC_CLIENT_ID=viz-ifc-consumer
VIZ_KAFKA_IFC_GROUP_ID=viz-ifc-consumers

KAFKA_BROKER=kafka:9093
KAFKA_IFC_TOPIC=ifc-files
```

Make sure to replace `ROOTUSER` and `CHANGEME123` with your own credentials.

## Running the services

To run the services, navigate to the root directory of the repository and run the following command:

```bash
docker-compose up
```

If you want to build the images before running the services, you can run the following command:

```bash
docker-compose up --build
```

## Stopping the services

To stop the services, navigate to the root directory of the repository and run the following command:

```bash
docker-compose down
```

## Access the MinIO Console

To access the MinIO Console, navigate to http://localhost:9001. You will need to use the credentials `ROOTUSER` and `CHANGEME123` to log in. Or whatever you have set in the `.env` file.

## Access MongoDB

To access the MongoDB database, navigate to http://localhost:27017.

## Tests

### Unit Tests

Unit tests are set up for every consumer service. They are running automatically when you make a pull request to main.

To run the unit tests locally, navigate to the service directory and run `npm test`.

### Integration Tests

Follow the instructions in the [integration-tests/README.md](integration-tests/README.md) file.
