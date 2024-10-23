import { Worker } from "worker_threads";
import path from "path";
import { WorkerResult } from "./ifcWorker";

/**
 * Runs a worker thread to process the given file buffer.
 *
 * @param file - The file buffer to be processed by the worker.
 * @returns  A promise that resolves with the result from the worker or rejects with an error.
 *
 * @throws If the worker stops with a non-zero exit code.
 */
function runWorker(file: Buffer, fileName: string): Promise<WorkerResult> {
  return new Promise((resolve, reject) => {
    const workerPath = path.resolve(__dirname, "ifcWorker.js");
    console.log(workerPath);
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

/**
 * Parses an IFC file and converts it to fragments.
 *
 * @param file - The IFC file to be parsed.
 * @param fileName - The name of the IFC file.
 * @returns A promise that resolves to the parsed fragments as a Buffer, or undefined if an error occurs.
 */
export default async function parseIfcToFragments(file: Buffer, fileName: string): Promise<Buffer> {
  const result = await runWorker(file, fileName);
  if (result.success) {
    if (!result.fragments) {
      console.error(`Error processing file ${fileName}: No fragments returned`);
      return Buffer.from([]);
    }
    console.log(`File ${fileName} processed successfully`);
    console.log(`Fragments size: ${result.fragments?.length} bytes`);
    console.log("result.Fragments is buffer:", Buffer.isBuffer(result.fragments));
    return result.fragments;
  } else {
    console.error(`Error processing file ${fileName}: ${result.error}`);
    return Buffer.from([]);
  }
}
