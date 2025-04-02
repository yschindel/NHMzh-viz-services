/**
 * IFC consumer module
 *
 * This module is the entry point for the IFC consumer.
 *
 * @module index
 */

import { runIfcToGzWorker } from "./ifc/ifcParser";
import { ensureWasmFile } from "./ifc/wasm";
import { setupKafkaConsumer, startKafkaConsumer } from "./kafka";
import { getFile, getFileMetadata, minioClient } from "./minio";
import { log } from "./utils/logger";
import { getEnv } from "./utils/env";

const IFC_BUCKET_NAME = getEnv("MINIO_IFC_BUCKET");

/**
 * Main function to start the IFC consumer
 * Ensures the WASM file is downloaded
 * Sets up the Kafka consumer
 * Starts the Kafka consumer
 */
async function main() {
	log.info("Starting server...");

	ensureWasmFile();

	log.info("Setting up Kafka consumer...");
	const consumer = await setupKafkaConsumer();
	log.info("Kafka consumer setup complete");

	log.info("Starting Kafka consumer...");
	await startKafkaConsumer(consumer, async (message: any) => {
		if (message.value) {
			try {
				const donwloadLink = message.value.toString();
				const fileID = donwloadLink.split("/").pop();

				const metadata = await getFileMetadata(fileID, IFC_BUCKET_NAME, minioClient);
				const ifcData: IFCData = {
					Project: metadata.project,
					Filename: metadata.filename,
					Timestamp: metadata.timestamp,
					File: await getFile(fileID, IFC_BUCKET_NAME, minioClient),
					FileID: fileID,
				};

				// Parse the IFC file to fragments and save to minio in a worker thread - no need to await
				runIfcToGzWorker(ifcData);
			} catch (error: any) {
				log.error("Error processing Kafka message:", error);
			}
		}
	});

	log.info("Kafka consumer started");
}

if (require.main === module) {
	main().catch(log.error);
}
