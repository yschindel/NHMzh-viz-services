import { setupKafkaConsumer, startKafkaConsumer } from "./kafka";
import { saveToMinIO } from "./minio";

async function main() {
  const consumer = await setupKafkaConsumer();

  await startKafkaConsumer(consumer, async (message) => {
    if (message.value) {
      const filename = message.key?.toString() || "unknown.ifc";
      await saveToMinIO(filename, message.value);
    }
  });

  console.log("IFC Consumer is running...");
}

if (require.main === module) {
  main().catch(console.error);
}
