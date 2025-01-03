# LCA and Cost Service

Responsibilites:

- Kafka consumer that listens to the LCA Kafka topic
- Save the LCA data as to Azure SQL DB (merges with Cost data if present)
- Kafka consumer that listens to the Cost Kafka topic
- Save the Cost data as Azure SQL DB (merges with LCA data if present)

## Kafka

The service has two Kafka consumers, one for the LCA data and one for the Cost data.

Each consumer subscribes to its own Kafka topic, which is specified in the `KAFKA_LCA_TOPIC` and `KAFKA_COST_TOPIC` environment variables.

## MinIO

The service has one bucket in MinIO:

- `lca-cost-data`: Where the LCA and Cost data is saved.

## Azure

Ensure one table in Azure SQL Server:

```sql
CREATE TABLE project_data (
    project NVARCHAR(255),
    filename NVARCHAR(255),
    timestamp NVARCHAR(255),
    id NVARCHAR(255),
    category NVARCHAR(255),
    material_kbob NVARCHAR(255),
    gwp_absolute FLOAT,
    gwp_relative FLOAT,
    penr_absolute FLOAT,
    penr_relative FLOAT,
    ubp_absolute FLOAT,
    ubp_relative FLOAT,
    cost FLOAT,
    cost_unit FLOAT
)
```
