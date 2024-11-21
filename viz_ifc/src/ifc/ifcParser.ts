import { Worker } from "worker_threads";
import path from "path";

/**
 * Runs a worker thread to process the given ifc file into a compressed 'fragments' file.
 * Saves the compressed file to MinIO.
 *
 * @param file - The file buffer to be processed by the worker.
 * @param fileName - The name of the file to be processed.
 * @throws If the worker stops with a non-zero exit code.
 */
export function runIfcToGzWorker(file: Buffer, fileName: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const workerPath = path.resolve(__dirname, "ifcWorker.js");
    const worker = new Worker(workerPath, {
      workerData: { file, fileName },
    });

    worker.on("message", resolve);
    worker.on("error", reject);
    worker.on("exit", (code) => {
      if (code !== 0) reject(new Error(`Worker stopped with exit code ${code}`));
    });
  });
}
