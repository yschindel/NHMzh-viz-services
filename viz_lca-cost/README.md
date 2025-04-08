# LCA and Cost Service

This services exists because:

The long-term data storage and reporting require a different data format than what is passed around in the Kafka network.

Responsibilites:

- Kafka consumer that listens to the LCA Kafka topic
- Kafka consumer that listens to the Cost Kafka topic
- Forward the processed data to another service

## Kafka

The service has two Kafka consumers, one for the LCA data and one for the Cost data.

Each consumer subscribes to its own Kafka topic, which is specified in the `KAFKA_LCA_TOPIC` and `KAFKA_COST_TOPIC` environment variables.

## Storage

The service will process the kafka messages and forward them as a POST request to a given REST API. This ensures that the core elements of the NHM systems which are closely integrated with Kafka and with one another are independend of the long-term storage solution required by the client for their reporting.
