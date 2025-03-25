import { Kafka, Consumer } from "kafkajs";

/**
 * Setup the Kafka consumer
 * @returns The Kafka consumer
 */
export async function setupKafkaConsumer(): Promise<Consumer> {
	const kafka = new Kafka({
		clientId: process.env.VIZ_KAFKA_IFC_CONSUMER_ID || "viz-ifc-consumer",
		brokers: [process.env.KAFKA_BROKER || "kafka:9093"],
	});

	const consumer = kafka.consumer({ groupId: process.env.VIZ_KAFKA_IFC_GROUP_ID || "viz-ifc-consumers" });

	try {
		await consumer.connect();
		await consumer.subscribe({ topic: process.env.KAFKA_IFC_TOPIC || "ifc-files", fromBeginning: true });
		return consumer;
	} catch (error) {
		console.error("Failed to connect to Kafka:", error);
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
				console.log(`Received message from topic ${topic}, partition ${partition}, offset ${message.offset}`);
				await messageHandler(message);
			},
		});
	} catch (error) {
		console.error("Error running Kafka consumer:", error);
		// Exit with a non-zero code to trigger restart
		process.exit(1);
	}
}
