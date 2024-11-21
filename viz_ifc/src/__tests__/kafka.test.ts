import { Kafka } from "kafkajs";
import { setupKafkaConsumer, startKafkaConsumer } from "../kafka";

// Mock KafkaJS
jest.mock("kafkajs");

const kafkaBroker = "test-broker:9092";
const kafkaClientId = "test-viz-ifc-consumer";
const kafkaGroupId = "test-viz-ifc-consumers";
const kafkaTopic = "test-ifc-files";

describe("Kafka Consumer", () => {
	let mockConsumer: any;
	let mockKafka: any;

	beforeEach(() => {
		mockConsumer = {
			connect: jest.fn(),
			subscribe: jest.fn(),
			run: jest.fn(),
		};

		mockKafka = {
			consumer: jest.fn(() => mockConsumer),
		};

		(Kafka as jest.MockedClass<typeof Kafka>).mockImplementation(() => mockKafka);

		// Mock environment variables
		process.env.KAFKA_BROKER = kafkaBroker;
		process.env.KAFKA_IFC_TOPIC = kafkaTopic;
		process.env.VIZ_KAFKA_IFC_GROUP_ID = kafkaGroupId;
		process.env.VIZ_KAFKA_IFC_CONSUMER_ID = kafkaClientId;
	});

	afterEach(() => {
		// Clean up environment variables after each test
		delete process.env.KAFKA_BROKER;
		delete process.env.KAFKA_IFC_TOPIC;
		delete process.env.VIZ_KAFKA_IFC_GROUP_ID;
		delete process.env.VIZ_KAFKA_IFC_CONSUMER_ID;
	});

	it("should create and connect a consumer", async () => {
		const consumer = await setupKafkaConsumer();

		expect(Kafka).toHaveBeenCalledWith({
			clientId: kafkaClientId,
			brokers: [kafkaBroker],
		});

		expect(mockKafka.consumer).toHaveBeenCalledWith({ groupId: kafkaGroupId });
		expect(mockConsumer.connect).toHaveBeenCalled();
		expect(mockConsumer.subscribe).toHaveBeenCalledWith({ topic: kafkaTopic, fromBeginning: true });

		expect(consumer).toBe(mockConsumer);
	});

	it("should run the consumer with the provided message handler", async () => {
		const messageHandler = jest.fn();
		await startKafkaConsumer(mockConsumer, messageHandler);

		expect(mockConsumer.run).toHaveBeenCalled();
		const runConfig = mockConsumer.run.mock.calls[0][0];
		expect(runConfig).toHaveProperty("eachMessage");

		// Simulate a message
		const message = { value: "test message" };
		await runConfig.eachMessage({ message });

		expect(messageHandler).toHaveBeenCalledWith(message);
	});

	it("should use default values when environment variables are not set", async () => {
		// Clear all environment variables
		delete process.env.KAFKA_BROKER;
		delete process.env.KAFKA_IFC_TOPIC;
		delete process.env.VIZ_KAFKA_IFC_GROUP_ID;
		delete process.env.VIZ_KAFKA_IFC_CONSUMER_ID;

		await setupKafkaConsumer();

		expect(Kafka).toHaveBeenCalledWith({
			clientId: "viz-ifc-consumer",
			brokers: ["kafka:9093"],
		});

		expect(mockKafka.consumer).toHaveBeenCalledWith({ groupId: "viz-ifc-consumers" });
		expect(mockConsumer.subscribe).toHaveBeenCalledWith({ topic: "ifc-files", fromBeginning: true });
	});

	it("should handle errors during consumer setup", async () => {
		mockConsumer.connect.mockRejectedValue(new Error("Connection failed"));

		await expect(setupKafkaConsumer()).rejects.toThrow("Connection failed");
	});

	it("should handle errors during message processing", async () => {
		const errorHandler = jest.fn();
		const messageHandler = jest.fn().mockRejectedValue(new Error("Processing failed"));

		mockConsumer.run.mockImplementation(async (config: any) => {
			try {
				await config.eachMessage({ message: {} });
			} catch (error) {
				errorHandler(error);
			}
		});

		await startKafkaConsumer(mockConsumer, messageHandler);

		expect(messageHandler).toHaveBeenCalled();
		expect(errorHandler).toHaveBeenCalledWith(expect.any(Error));
	});

	it("should handle different message formats", async () => {
		const messageHandler = jest.fn();
		await startKafkaConsumer(mockConsumer, messageHandler);

		const runConfig = mockConsumer.run.mock.calls[0][0];

		// Test with a message containing a Buffer
		const bufferMessage = { value: Buffer.from("test buffer message") };
		await runConfig.eachMessage({ message: bufferMessage });
		expect(messageHandler).toHaveBeenCalledWith(bufferMessage);

		// Test with a message containing a string
		const stringMessage = { value: "test string message" };
		await runConfig.eachMessage({ message: stringMessage });
		expect(messageHandler).toHaveBeenCalledWith(stringMessage);

		// Test with a message containing JSON
		const jsonMessage = { value: JSON.stringify({ key: "value" }) };
		await runConfig.eachMessage({ message: jsonMessage });
		expect(messageHandler).toHaveBeenCalledWith(jsonMessage);
	});
});
