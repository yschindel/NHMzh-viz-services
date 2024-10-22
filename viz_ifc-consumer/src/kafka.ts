import { Kafka, Consumer } from "kafkajs";

export async function setupKafkaConsumer(): Promise<Consumer> {
  const kafka = new Kafka({
    clientId: "ifc-consumer",
    brokers: [process.env.KAFKA_BROKER || "kafka:9092"],
  });

  const consumer = kafka.consumer({ groupId: "ifc-consumer-group" });
  await consumer.connect();
  await consumer.subscribe({ topic: "ifc-files", fromBeginning: true });

  return consumer;
}

export async function startKafkaConsumer(consumer: Consumer, messageHandler: (message: any) => Promise<void>): Promise<void> {
  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      await messageHandler(message);
    },
  });
}
