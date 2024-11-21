import { Kafka, Consumer, Producer } from "kafkajs";

/**
 * TODO: CONVERT TO RDF INPUT?? probably not, because we can't have buffers in RDF
 * Interface for the Kafka message payload
 * This has to be in sync with the producer
 * It's important to have the timestamp in the message payload to keep track of the file version
 * The project is also relevant, because file names can be the same across multiple projects.
 * The location is the path inside the MinIO bucket, consisting of project/filename_timestamp.extension
 * The locaion is used to retrieve the IFC file from the IFC MinIO Instance
 */
export interface IfcMsg {
	project: string;
	filename: string;
	timestamp: string; // Timestamp for the filename (in ISO 8601 format)
	location: string; // The file path inside the MinIO bucket, consisting of project/filename_timestamp
}

export function newIfcMsg(project: string, filename: string, timestamp: string, location: string): IfcMsg {
	return { project, filename, timestamp, location };
}

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
export async function sendKafkaMessage(producer: Producer, message: IfcMsg): Promise<void> {
	await producer.send({
		topic: process.env.KAFKA_IFC_TOPIC || "ifc-files",
		messages: [{ value: JSON.stringify(message) }],
	});
}
