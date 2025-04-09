/**
 * Kafka consumer integration module
 *
 * This module provides functionality for setting up and starting a Kafka consumer.
 *
 * @module kafka
 */

import { Kafka, Consumer } from "kafkajs";
import { getEnv } from "./utils/env";
import { log } from "./utils/logger";

/**
 * Setup the Kafka consumer
 * @returns The Kafka consumer
 */
export async function setupKafkaConsumer(): Promise<Consumer> {
	const kafka = new Kafka({
		clientId: "viz-ifc-consumer",
		brokers: [getEnv("KAFKA_BROKER")],

		// Custom log creator to log Kafka messages to align with the logger in this service
		logCreator:
			() =>
			({ namespace, level, label, log }) => {
				const logLevel = String(level).toLowerCase();
				log.log = { namespace, label, ...log };
				log.logger?.includes("kafkajs") && log.message && log.log;
				if (logLevel === "error") {
					log.error(log.message, log.log);
				} else if (logLevel === "warn") {
					log.warn(log.message, log.log);
				} else if (logLevel === "info") {
					log.info(log.message, log.log);
				} else if (logLevel === "debug") {
					log.debug(log.message, log.log);
				}
			},
	});

	const consumer = kafka.consumer({
		groupId: "viz-ifc-consumer-group",
		// Add session timeout for long-running processes
		sessionTimeout: 120000, // 2 minutes
	});

	try {
		log.debug("Connecting to Kafka");
		await consumer.connect();
		log.debug("Subscribing to Kafka topic");
		// Change fromBeginning to false to only process new messages
		await consumer.subscribe({
			topic: getEnv("KAFKA_IFC_TOPIC"),
			fromBeginning: false,
		});
		log.debug("Kafka consumer connected and subscribed to topic");
		return consumer;
	} catch (error: any) {
		log.error("Failed to connect to Kafka:", error);
		// Exit with a non-zero code to trigger restart
		process.exit(1);
	}
}

/**
 * Start the Kafka consumer
 * @param consumer - The Kafka consumer
 * @param messageHandler - The message handler
 */
export async function startKafkaConsumer(consumer: Consumer, messageHandler: (message: any) => Promise<void>): Promise<void> {
	try {
		await consumer.run({
			autoCommit: false, // Disable auto-commit
			eachMessage: async ({ topic, partition, message, heartbeat }) => {
				try {
					log.debug(`Processing message from topic ${topic}, partition ${partition}, offset ${message.offset}`);
					await messageHandler(message);

					// Only commit after successful processing
					await consumer.commitOffsets([
						{
							topic,
							partition,
							offset: (Number(message.offset) + 1).toString(),
						},
					]);
					log.debug(`Committed offset ${Number(message.offset) + 1} for topic ${topic}, partition ${partition}`);
				} catch (error: any) {
					log.error("Error processing message:", error);
					// Don't commit the offset if processing failed
					// The message will be reprocessed
				}
			},
		});
	} catch (error: any) {
		log.error("Error running Kafka consumer:", error);
		// Exit with a non-zero code to trigger restart
		process.exit(1);
	}
}
