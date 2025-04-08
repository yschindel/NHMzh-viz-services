/**
 * Kafka consumer integration module
 *
 * This module provides functionality for setting up and starting a Kafka consumer.
 *
 * @module kafka
 */

import { Kafka } from "kafkajs";
import type { Consumer, Producer } from "kafkajs";
import { getEnv } from "../utils/env";
import { log } from "../utils/logger";
import type { IfcFileData } from "../types";

/**
 * Setup the Kafka consumer
 * @returns The Kafka consumer
 */
export async function setupKafkaConsumer(): Promise<Consumer> {
	log.info("Setting up Kafka consumer...");
	const kafka = new Kafka({
		clientId: "mock_calc_lca-cost-consumer",
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

	const consumer = kafka.consumer({ groupId: "mock_calc_lca-cost-consumer-group" });

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
	log.info("Starting Kafka consumer...");
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

/**
 * Setup the Kafka producer
 * @returns The Kafka producer
 */
export async function setupKafkaProducer(): Promise<Producer> {
	log.info("Setting up Kafka producer...");
	const kafka = new Kafka({
		clientId: "mock_calc_lca-cost-producer",
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
	return kafka.producer();
}

/**
 * Start the Kafka producer
 * @param producer - The Kafka producer
 */
export async function startKafkaProducer(producer: Producer): Promise<void> {
	log.info("Starting Kafka producer...");
	try {
		await producer.connect();
		log.debug("Kafka producer connected");
	} catch (error: any) {
		log.error("Failed to connect to Kafka:", error);
		// Exit with a non-zero code to trigger restart
		process.exit(1);
	}
}

/**
 * Send a message to the Kafka producer
 * @param producer - The Kafka producer
 * @param message - The message to send
 */
export async function sendMessage(producer: Producer, message: IfcFileData, topic: string): Promise<void> {
	try {
		await producer.send({
			topic,
			messages: [{ value: JSON.stringify(message) }],
		});
	} catch (error: any) {
		log.error("Failed to send message to Kafka:", error);
	}
}
