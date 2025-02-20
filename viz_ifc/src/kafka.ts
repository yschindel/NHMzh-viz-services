import { Kafka, Consumer, Producer } from "kafkajs";

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
	await consumer.connect();
	await consumer.subscribe({ topic: process.env.KAFKA_IFC_TOPIC || "ifc-files", fromBeginning: true });

	return consumer;
}

/**
 * Start the Kafka consumer
 * @param consumer - The Kafka consumer
 * @param messageHandler - The message handler
 */
export async function startKafkaConsumer(consumer: Consumer, messageHandler: (message: any) => Promise<void>): Promise<void> {
	await consumer.run({
		eachMessage: async ({ topic, partition, message }) => {
			console.log(`Received message from topic ${topic}, partition ${partition}, offset ${message.offset}`);
			await messageHandler(message);
		},
	});
}

/**
 * Setup the Kafka producer
 * @returns The Kafka producer
 */
export async function setupKafkaProducer(): Promise<Producer> {
	const kafka = new Kafka({
		clientId: process.env.VIZ_KAFKA_IFC_PRODUCER_ID || "viz-ifc-producer",
		brokers: [process.env.KAFKA_BROKER || "kafka:9093"],
	});

	const producer = kafka.producer();
	await producer.connect();
	return producer;
}

/**
 * Send a message to the Kafka producer
 * @param producer - The Kafka producer
 * @param message - The message to send
 */
export async function sendKafkaMessage(producer: Producer, message: string): Promise<void> {
	await producer.send({
		topic: process.env.KAFKA_IFC_TOPIC || "ifc-files",
		messages: [{ value: message }],
	});
}
