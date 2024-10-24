import { Kafka, logLevel } from "kafkajs";
import * as path from "path";
import * as dotenv from "dotenv";

dotenv.config({ path: path.resolve(__dirname, "../.env") });

const KAFKA_BROKER = "localhost:9092";
const TOPIC = process.env.KAFKA_LCA_TOPIC || "lca-data";

const MESSAGES = [createMessage("project1", "file1.ifc"), createMessage("project2", "file2.ifc"), createMessage("project1", "file3.ifc")];

interface Message {
  project: string;
  filename: string;
  data: Element[]
  timestamp: string;
}

interface Element {
  category: string;
  "CO2-eq (kg)": number;
  "Graue Energie": number;
  UBP: number;
}

function createMessage(project: string, filename: string): Message {
  const timestamp = new Date().toISOString();
  return {
    project,
    filename,
    data: createData(),
    timestamp,
  };
}

export function createData(): any[] {
  const numElements = Math.floor(Math.random() * (10000 - 1000 + 1)) + 1000;
  const elements = [];
  for (let i = 0; i < numElements; i++) {
    elements.push(createElementData());
  }
  return elements;
}

function createElementData() {
  return {
    "category": getRandomCategory(),
    "CO2-eq (kg)": randomInt(),
    "Graue Energie": randomInt(),
    "UBP": randomInt(),
  };
}

function randomInt() {
  return Math.floor(Math.random() * 100);
}

function getRandomCategory() {
  const categories = ["Wall", "Door", "Floor"];
  return categories[Math.floor(Math.random() * categories.length)];
}

// Kafka Stuff

async function produceMessages(): Promise<void> {
  const kafka = new Kafka({
    clientId: "my-app",
    brokers: [KAFKA_BROKER],
    logLevel: logLevel.ERROR,
  });

  const producer = kafka.producer();
  await producer.connect();

  for (const msg of MESSAGES) {
    await producer.send({
      topic: TOPIC,
      messages: [{ value: JSON.stringify(msg) }],
    });
    console.log(`Sent message: ${msg.project}/${msg.filename}`);
  }

  await producer.disconnect();
}

// consuming messages from Kafka
// just for validating that the messages are being sent
// actual consumption of messages is done in the viz_lca-consumer
async function consumeMessages(): Promise<void> {
  const kafka = new Kafka({
    clientId: "my-app",
    brokers: [KAFKA_BROKER],
    logLevel: logLevel.ERROR,
  });

  const consumer = kafka.consumer({ groupId: "test-group" });
  await consumer.connect();
  await consumer.subscribe({ topic: TOPIC, fromBeginning: true });

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      if (!message || !message.value) {
        return;
      }
      const value = JSON.parse(message.value.toString());
      console.log(`Received message: ${value.project}/${value.filename}`);
      console.log(`Content: ${value.location}`);
    },
  });
}

(async () => {
  console.log("Producing LCA messages...");
  await produceMessages();
  console.log("\nConsuming LCA messages...");
  await consumeMessages();
})();
