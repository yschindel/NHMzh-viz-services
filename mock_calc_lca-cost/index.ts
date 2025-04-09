/**
 * IFC consumer module
 *
 * This module is the entry point for the IFC consumer.
 *
 * @module index
 */

import { sendMessage, setupKafkaConsumer, startKafkaConsumer, setupKafkaProducer, startKafkaProducer } from "./src/services/kafka";
import { getFile, getFileMetadata, minioClient } from "./src/services/minio";
import { log } from "./src/utils/logger";
import { getEnv } from "./src/utils/env";

import { processIfc } from "./src/services/ifcProcessor";
import type { IfcFileData } from "./src/types";
import { chunkArray } from "./src/utils/chunk";

const IFC_BUCKET_NAME = getEnv("MINIO_IFC_BUCKET");
const LCA_TOPIC = getEnv("KAFKA_LCA_TOPIC");
const COST_TOPIC = getEnv("KAFKA_COST_TOPIC");

/**
 * Main function to start the IFC consumer
 * Ensures the WASM file is downloaded
 * Sets up the Kafka consumer
 * Starts the Kafka consumer
 * Sets up the Kafka producer
 * Starts the Kafka producer
 */
async function main() {
	log.info("Starting server...");

	const consumer = await setupKafkaConsumer();
	const producer = await setupKafkaProducer();
	await startKafkaProducer(producer);

	log.info("Starting Kafka consumer...");
	await startKafkaConsumer(consumer, async (message: any) => {
		if (message.value) {
			try {
				log.info("Processing Kafka message:", message.value);
				const downloadLink = message.value.toString();
				const fileID = downloadLink.split("/").pop();
				const file = await getFile(fileID, IFC_BUCKET_NAME, minioClient);
				if (!file) {
					log.error(`File ${fileID} not found`);
					return;
				}

				log.debug("Getting file metadata...");
				const metadata = await getFileMetadata(fileID, IFC_BUCKET_NAME, minioClient);

				log.debug("Processing IFC file...");
				const { lcaData, costData } = await processIfc(file);

				const ifcFileData: IfcFileData = {
					project: metadata.project,
					filename: metadata.filename,
					timestamp: metadata.timestamp,
				};

				// chunk the lcaData and costData array into 1000 element chunks
				const lcaDataChunks = chunkArray(lcaData, 1000);
				const costDataChunks = chunkArray(costData, 1000);

				log.debug("Sending LCA data to Kafka...");
				if (lcaData.length > 0) {
					for (const chunk of lcaDataChunks) {
						ifcFileData.data = chunk;
						await sendMessage(producer, ifcFileData, LCA_TOPIC);
					}
				}

				log.debug("Sending cost data to Kafka...");
				if (costData.length > 0) {
					for (const chunk of costDataChunks) {
						ifcFileData.data = chunk;
						await sendMessage(producer, ifcFileData, COST_TOPIC);
					}
				}
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
