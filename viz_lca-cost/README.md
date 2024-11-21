# LCA and Cost Service

Responsibilites:

- Kafka consumer that listens to the LCA Kafka topic
- Save the LCA data as .parquet files to MinIO (merges with Cost data if present)
- Kafka consumer that listens to the Cost Kafka topic
- Save the Cost data as .parquet files to MinIO (merges with LCA data if present)

## Kafka

The service has two Kafka consumers, one for the LCA data and one for the Cost data.

Each consumer subscribes to its own Kafka topic, which is specified in the `KAFKA_LCA_TOPIC` and `KAFKA_COST_TOPIC` environment variables.

## DuckDB

The service uses DuckDB to merge the LCA and Cost data into a single .parquet file.

## MinIO

The service has one bucket in MinIO:

- `lca-cost-data`: Where the LCA and Cost data is saved.
