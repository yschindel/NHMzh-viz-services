import { setupKafkaConsumer, startKafkaConsumer, IFCKafkaMessage } from "./kafka";
import { initializeMinio, saveToMinIO } from "./minio";

const BUCKET_NAME = process.env.MINIO_IFC_FRAGMENTS_BUCKET || "ifc-fragment-files";

/**
 * Main function to start the IFC consumer
 * Initialize the Minio bucket and start the Kafka consumer
 */
async function main() {
  await initializeMinio(BUCKET_NAME);
  const consumer = await setupKafkaConsumer();

  await startKafkaConsumer(consumer, async (message: any) => {
    if (message.value) {
      try {
        const ifcMessage: IFCKafkaMessage = JSON.parse(message.value.toString());
        await saveToMinIO(ifcMessage.project, ifcMessage.filename, Buffer.from(ifcMessage.content), BUCKET_NAME, ifcMessage.timestamp);
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
