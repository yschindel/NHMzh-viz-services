import parseIfcToFragments from "./ifc/ifcParser";
import { setupKafkaConsumer, startKafkaConsumer, IFCKafkaMessage } from "./kafka";
import { initializeMinio, saveToMinIO, getFile, createMinioClient } from "./minio";

const FRAGMENTS_BUCKET_NAME = process.env.MINIO_IFC_FRAGMENTS_BUCKET || "ifc-fragment-files";
const IFC_BUCKET_NAME = process.env.MINIO_IFC_BUCKET || "ifc-files";

/**
 * Main function to start the IFC consumer
 * Initialize the Minio bucket and start the Kafka consumer
 */
async function main() {
  const minioClient = createMinioClient();
  await initializeMinio(FRAGMENTS_BUCKET_NAME, minioClient);
  const consumer = await setupKafkaConsumer();

  await startKafkaConsumer(consumer, async (message: any) => {
    if (message.value) {
      try {
        const value = JSON.parse(message.value.toString());
        const ifcMessage: IFCKafkaMessage = {
          project: value.project,
          filename: value.filename,
          location: value.location,
          timestamp: value.timestamp,
        };

        const file = await getFile(ifcMessage.location, IFC_BUCKET_NAME, minioClient);
        const fragments = await parseIfcToFragments(file, ifcMessage.location);
        if (fragments) {
          console.log("fragments is buffer:", Buffer.isBuffer(fragments));

          await saveToMinIO(minioClient, FRAGMENTS_BUCKET_NAME, ifcMessage.project, ifcMessage.filename, ifcMessage.timestamp, fragments);
        }
      } catch (error) {
        console.error("Error processing Kafka message:", error);
      }
    }
  });

  console.log("IFC Consumer is running...");
}

if (require.main === module) {
  main().catch(console.error);
}
