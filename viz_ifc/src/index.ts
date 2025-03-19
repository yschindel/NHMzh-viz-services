import { runIfcToGzWorker } from "./ifc/ifcParser";
import { ensureWasmFile } from "./ifc/wasm";
import { setupKafkaConsumer, startKafkaConsumer } from "./kafka";
import { initializeMinio, getFile, getFileMetadata, minioClient } from "./minio";

const FRAGMENTS_BUCKET_NAME = process.env.MINIO_IFC_FRAGMENTS_BUCKET || "ifc-fragment-files";
const IFC_BUCKET_NAME = process.env.MINIO_IFC_BUCKET || "ifc-files";

/**
 * Main function to start the IFC consumer
 * Initialize the Minio bucket and start the Kafka consumer
 */
async function main() {
	ensureWasmFile();
	await initializeMinio([FRAGMENTS_BUCKET_NAME, IFC_BUCKET_NAME], minioClient);

	const consumer = await setupKafkaConsumer();
	console.log("IFC Consumer is running...");

	await startKafkaConsumer(consumer, async (message: any) => {
		if (message.value) {
			try {
				console.log("received message: ", message);
				const donwloadLink = message.value.toString();
				const location = donwloadLink.split("/").pop();

				console.log("getting object from minio at: ", location);
				const file = await getFile(location, IFC_BUCKET_NAME, minioClient);

				console.log("getting metadata from minio at: ", location);
				const metadata = await getFileMetadata(location, IFC_BUCKET_NAME, minioClient);

				console.log("received metadata: ", metadata);

				// Parse the IFC file to fragments and save to minio in a worker thread - no need to await
				console.log("running ifc to gz worker");
				runIfcToGzWorker(file, location, metadata.timestamp, metadata.project, metadata.filename);
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
