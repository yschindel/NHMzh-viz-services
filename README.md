# NHMzh Vizualization Services

This repository contains the code for the Vizualization Services of the NHMzh project.

## Prerequisites

- Docker
- Docker Compose

## Environment Variables

Create a `.env` file in the root directory of the repository with the following variables:

```
MINIO_ROOT_USER=ROOTUSER
MINIO_ROOT_PASSWORD=CHANGEME123
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

To access the MinIO Console, navigate to http://localhost:9001. You will need to use the credentials `ROOTUSER` and `CHANGEME123` to log in.

## Access MongoDB

To access the MongoDB database, navigate to http://localhost:27017.
