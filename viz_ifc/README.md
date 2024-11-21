# IFC Service

Responsibilites:

- Api endpoint for uploading IFC files
- Save the IFC file to MinIO
- Kafka producer for notifications when a file is uploaded
- Kafka consumer for notifications when a file is uploaded
- Convert the IFC file to fragments and save them to MinIO

## Kafka

The Services is the producer and consumer for the IFC Kafka topic at the same time. This allows other services to hook into the event of a new IFC file being uploaded.

The service subscribes to a Kafka topic, which is specified in the `KAFKA_IFC_TOPIC` environment variable.

The Kafka messages have the following format:

```json
{
	"project": "string",
	"filename": "string",
	"timestamp": "number",
	"location": "string"
}
```

> Refer to the [IFCKafkaMessage interface](./src/kafka.ts) for more information.

## MinIO

The service has two buckets in MinIO:

- `ifc-files`: Where the original IFC files are saved.
- `ifc-fragment-files`: Where the fragments are saved.

The IFC upload API saves the files to the `ifc-files` bucket.
The converted fragments are saved to the `ifc-fragment-files` bucket.

The files are saved with a unique name, which is a combination of the project name, the original filename, and a timestamp.

Example: `project2/file2_2024-10-25T16:36:04.986173Z.gz`

## Fragmentation

The IFC files are converted to fragments using [@ThatOpen/](https://www.npmjs.com/org/thatopen)
This enables their use in the web viewer also provided by ThatOpen.
