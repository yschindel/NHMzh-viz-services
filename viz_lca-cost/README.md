# LCA and Cost Service

Responsibilites:

- Kafka consumer that listens to the LCA Kafka topic
- Save the LCA data as to Azure SQL DB
- Kafka consumer that listens to the Cost Kafka topic
- Save the Cost data as Azure SQL DB

## Kafka

The service has two Kafka consumers, one for the LCA data and one for the Cost data.

Each consumer subscribes to its own Kafka topic, which is specified in the `KAFKA_LCA_TOPIC` and `KAFKA_COST_TOPIC` environment variables.

## Azure

The service will create and ensure columns in an Azure SQL DB on startup. Database configuration is done in ./sql/
The following environment variables need to be set:

```
AZURE_DB_SERVER
AZURE_DB_PORT
AZURE_DB_USER
AZURE_DB_PASSWORD
AZURE_DB_DATABASE
```
