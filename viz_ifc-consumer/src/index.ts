import { Kafka } from "kafkajs";
import { Client } from "minio";
import * as fs from "fs";
import * as path from "path";

// Configure Kafka consumer
const kafka = new Kafka({
  clientId: "ifc-consumer",
  brokers: ["localhost:9092"], // Adjust if your Kafka broker is elsewhere
});

const consumer = kafka.consumer({ groupId: "ifc-consumer-group" });

// Configure MinIO client
const minioClient = new Client({
  endPoint: "localhost",
  port: 9000,
  useSSL: false,
  accessKey: "your-access-key", // Replace with your MinIO access key
  secretKey: "your-secret-key", // Replace with your MinIO secret key
});

async function saveToMinIO(fileName: string, fileContent: Buffer): Promise<void> {
  const bucketName = "ifc-files"; // Replace with your desired bucket name

  // Ensure the bucket exists
  const bucketExists = await minioClient.bucketExists(bucketName);
  if (!bucketExists) {
    await minioClient.makeBucket(bucketName);
  }

  // Save the file to MinIO
  await minioClient.putObject(bucketName, fileName, fileContent);
  console.log(`File ${fileName} saved to MinIO bucket ${bucketName}`);
}

async function run(): Promise<void> {
  await consumer.connect();
  await consumer.subscribe({ topic: "ifc-file", fromBeginning: true });

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      if (message.value) {
        const fileName = message.key?.toString() || `ifc-file-${Date.now()}.ifc`;
        const fileContent = message.value;

        console.log(`Received file: ${fileName}`);
        await saveToMinIO(fileName, fileContent);
      }
    },
  });
}

run().catch(console.error);
