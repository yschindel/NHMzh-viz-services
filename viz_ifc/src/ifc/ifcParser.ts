/**
 * IFC parser module
 *
 * This module provides functionality for running Web Workers to parse IFC files into fragments.
 *
 * @module ifcParser
 */

import { Worker } from "worker_threads";
import path from "path";
import { log } from "../utils/logger";

/**
 * Runs a worker thread to process the given ifc file into a compressed 'fragments' file.
 * Saves the compressed file to MinIO.
 *
 * @param file - The file buffer to be processed by the worker.
 * @param location - The location of the file in the bucket.
 * @param timestamp - The timestamp of the file.
 * @param project - The project name.
 * @throws If the worker stops with a non-zero exit code.
 */
export function runIfcToGzWorker(file: Buffer, location: string, timestamp: string, project: string, filename: string): Promise<void> {
	log.debug(`Running ifc to gz worker for file ${filename}`);
	return new Promise((resolve, reject) => {
		const workerPath = path.resolve(__dirname, "ifcWorker.js");
		const worker = new Worker(workerPath, {
			workerData: { file, location, timestamp, project, filename },
		});

		worker.on("message", resolve);
		worker.on("error", reject);
		worker.on("exit", (code) => {
			if (code !== 0) reject(new Error(`Worker stopped with exit code ${code}`));
		});
	});
}
