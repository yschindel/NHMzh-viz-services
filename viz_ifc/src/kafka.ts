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
		clientId: getEnv("VIZ_KAFKA_IFC_CONSUMER_ID"),
		brokers: [getEnv("KAFKA_BROKER")],
	});

	const consumer = kafka.consumer({ groupId: getEnv("VIZ_KAFKA_IFC_GROUP_ID") });

	try {
		log.debug("Connecting to Kafka");
		await consumer.connect();
		log.debug("Subscribing to Kafka topic");
		await consumer.subscribe({ topic: getEnv("KAFKA_IFC_TOPIC"), fromBeginning: true });
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
			eachMessage: async ({ topic, partition, message }) => {
				log.debug(`Received message from topic ${topic}, partition ${partition}, offset ${message.offset}`);
				await messageHandler(message);
			},
		});
	} catch (error: any) {
		log.error("Error running Kafka consumer:", error);
		// Exit with a non-zero code to trigger restart
		process.exit(1);
	}
}
