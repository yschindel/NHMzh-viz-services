import { runIfcToGzWorker } from "./ifc/ifcParser";
import { ensureWasmFile } from "./ifc/wasm";
import { setupKafkaConsumer, startKafkaConsumer, IfcMsg, setupKafkaProducer } from "./kafka";
import { initializeMinio, getFile, getFileMetadata, minioClient } from "./minio";
import express from "express";
import apiRouter from "./api";
import { Producer } from "kafkajs";

const FRAGMENTS_BUCKET_NAME = process.env.MINIO_IFC_FRAGMENTS_BUCKET || "ifc-fragment-files";
const IFC_BUCKET_NAME = process.env.MINIO_IFC_BUCKET || "ifc-files";
const IFC_API_PORT = process.env.IFC_API_PORT || 4242;

export let producer: Producer;

/**
 * Main function to start the IFC consumer
 * Initialize the Minio bucket and start the Kafka consumer
 */
async function main() {
	ensureWasmFile();
	await initializeMinio([FRAGMENTS_BUCKET_NAME, IFC_BUCKET_NAME], minioClient);

	const consumer = await setupKafkaConsumer();
	console.log("IFC Consumer is running...");

	producer = await setupKafkaProducer();
	console.log("IFC Producer is running...");

	// Setup Express API server
	const app = express();
	app.use("/api", apiRouter);

	app.listen(IFC_API_PORT, () => {
		console.log(`API server listening on port ${IFC_API_PORT}`);
	});

	await startKafkaConsumer(consumer, async (message: any) => {
		if (message.value) {
			try {
				const value = JSON.parse(message.value.toString());
				const ifcMessage: IfcMsg = {
					project: value.project,
					filename: value.filename,
					location: value.location,
					timestamp: value.timestamp,
				};

				const file = await getFile(ifcMessage.location, IFC_BUCKET_NAME, minioClient);
				const metadata = await getFileMetadata(ifcMessage.location, IFC_BUCKET_NAME, minioClient);
				console.log("metadata", metadata);

				// Parse the IFC file to fragments and save to minio in a worker thread - no need to await
				runIfcToGzWorker(file, ifcMessage.location, metadata.timestamp, metadata.project, metadata.filename);
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
