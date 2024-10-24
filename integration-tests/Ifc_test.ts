import { Kafka, logLevel } from "kafkajs";
import { Client as MinioClient } from "minio";
import * as fs from "fs";
import * as path from "path";
import * as dotenv from "dotenv";

dotenv.config({ path: path.resolve(__dirname, "../.env") });

const KAFKA_BROKER = "localhost:9092";
const TOPIC = process.env.KAFKA_IFC_TOPIC || "ifc-files";
const IFC_BUCKET = process.env.MINIO_IFC_BUCKET || "ifc-files";
const MINIO_ENDPOINT = "localhost";
const MINIO_PORT = parseInt(process.env.MINIO_PORT || "9000");
const MINIO_USE_SSL = process.env.MINIO_USE_SSL === "true";
const MINIO_ACCESS_KEY = process.env.MINIO_ACCESS_KEY || "";
const MINIO_SECRET_KEY = process.env.MINIO_SECRET_KEY || "";

const MESSAGES = [createMessage("project1", "file1.ifc"), createMessage("project2", "file2.ifc"), createMessage("project1", "file3.ifc")];

interface Message {
  project: string;
  filename: string;
  location: string;
  timestamp: string;
}

function createMessage(project: string, filename: string): Message {
  const timestamp = new Date().toISOString();
  return {
    project,
    filename,
    location: createFileName(project, filename, timestamp),
    timestamp,
  };
}


export function createFileName(project: string, filename: string, timestamp: string): string {
  const { name, ext } = path.parse(filename);
  const fileTimestamp = timestamp || new Date().toISOString();
  return `${project}/${name}_${fileTimestamp}${ext}`;
}

// MinIO Stuff

async function addIfcFileToMinio(client: MinioClient, location: string, content: Buffer): Promise<void> {
  await client.putObject(IFC_BUCKET, location, content);
}

async function addIfcFilesToMinio(): Promise<void> {
  const client = new MinioClient({
    endPoint: MINIO_ENDPOINT,
    port: MINIO_PORT,
    useSSL: MINIO_USE_SSL,
    accessKey: MINIO_ACCESS_KEY,
    secretKey: MINIO_SECRET_KEY,
  });

  await createBucket(client, IFC_BUCKET);
  const content = fs.readFileSync(path.join(__dirname, "assets/test.ifc"));
  for (const msg of MESSAGES) {
    await addIfcFileToMinio(client, msg.location, content);
  }
}

// create a bucket if it doesn't exist
async function createBucket(client: MinioClient, bucketName: string): Promise<void> {
  const exists = await client.bucketExists(bucketName);
  if (!exists) {
    await client.makeBucket(bucketName);
  }
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
// actual consumption of messages is done in the viz_ifc-consumer
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
  console.log("Adding IFC files to MinIO...");
  await addIfcFilesToMinio();
  console.log("Waiting for the files to be processed...");
  await new Promise((resolve) => setTimeout(resolve, 3000));
  console.log("Producing IFC messages...");
  await produceMessages();
  console.log("\nConsuming IFC messages...");
  await consumeMessages();
})();
