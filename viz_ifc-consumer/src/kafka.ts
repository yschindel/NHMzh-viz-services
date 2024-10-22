import { Kafka, Consumer } from "kafkajs";

export async function setupKafkaConsumer(): Promise<Consumer> {
  const kafka = new Kafka({
    clientId: process.env.VIZ_KAFKA_IFC_CLIENT_ID || "viz-ifc-consumer",
    brokers: [process.env.KAFKA_BROKER || "kafka:9092"],
  });

  const consumer = kafka.consumer({ groupId: process.env.VIZ_KAFKA_IFC_GROUP_ID || "viz-ifc-consumers" });
  await consumer.connect();
  await consumer.subscribe({ topic: process.env.KAFKA_IFC_TOPIC || "ifc-files", fromBeginning: true });

  return consumer;
}

export async function startKafkaConsumer(consumer: Consumer, messageHandler: (message: any) => Promise<void>): Promise<void> {
  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      await messageHandler(message);
    },
  });
}
