import { Worker } from "worker_threads";
import path from "path";
import { IFCKafkaMessage } from "../kafka";
import { WorkerResult } from "./ifcWorker";

const __dirname = path.dirname(__filename);

function runWorker(file: Buffer): Promise<WorkerResult> {
  return new Promise((resolve, reject) => {
    const workerPath = path.resolve(__dirname, "ifcWorker.js");
    console.log(workerPath);
    const worker = new Worker(workerPath, {
      workerData: { file },
    });

    worker.on("message", resolve);
    worker.on("error", reject);
    worker.on("exit", (code) => {
      if (code !== 0) reject(new Error(`Worker stopped with exit code ${code}`));
    });
  });
}

export default async function parseIfc(kafkaMessage: IFCKafkaMessage): Promise<Buffer | undefined> {
  const result = await runWorker(kafkaMessage.content);
  if (result.success) {
    console.log(`File ${kafkaMessage.filename} processed successfully`);
    return result.fragments;
  } else {
    console.error(`Error processing file ${kafkaMessage.filename}: ${result.error}`);
  }
}
