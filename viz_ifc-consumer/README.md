# IFC Consumer

This service consumes IFC files from a Kafka topic and saves a compressed version of them to MinIO.

## Kafka

The service subscribes to a Kafka topic, which is specified in the `KAFKA_IFC_TOPIC` environment variable.

The Kafka messages have the following format:

```json
{
  "project": "string",
  "filename": "string",
  "timestamp": "number",
  "content": "Buffer"
}
```

> Refer to the [IFCKafkaMessage interface](./src/kafka.ts) for more information.

## MinIO

The service saves the IFC files to MinIO, in a bucket specified in the `MINIO_IFC_FRAGMENTS_BUCKET` environment variable.

The files are saved with a unique name, which is a combination of the project name, the original filename, and a timestamp.

Example: `my-project/my-file_1724230400000.gz`

## Fragmentation

The IFC files are converted to fragments using [@ThatOpen/](https://www.npmjs.com/org/thatopen)
This enables their use in the web viewer also provided by ThatOpen.
