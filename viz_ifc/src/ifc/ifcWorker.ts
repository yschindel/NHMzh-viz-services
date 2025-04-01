/**
 * IFC worker module
 *
 * This module provides functionality for parsing IFC files into fragments using Web Workers.
 *
 * @module ifcWorker
 */

import { parentPort, workerData } from "worker_threads";
import pako from "pako";
import * as OBC from "@thatopen/components";
import { wasmDir } from "./wasm";
import { sendFileToStorage } from "../storage";
import { log } from "../utils/logger";
import { getEnv } from "../utils/env";

const FRAGMENTS_BUCKET_NAME = getEnv("VIZ_IFC_FRAGMENTS_BUCKET");

export interface WorkerResult {
	success: boolean;
	fragments?: Buffer;
	error?: string;
}

/**
 * Converts IFC file to fragments and returns compressed data.
 * @param file The IFC file
 * @returns The result of the worker
 */
async function ifcToFragments(file: Buffer, fileName: string): Promise<WorkerResult> {
	log.info(`Converting IFC file ${fileName} to fragments`);
	const dataArray = new Uint8Array(file);
	const components = new OBC.Components();
	const fragments = components.get(OBC.FragmentsManager);
	const loader = components.get(OBC.IfcLoader);

	loader.settings.wasm = {
		path: wasmDir,
		absolute: true,
	};

	log.debug(`Loader setup with WASM path: ${wasmDir}`);

	log.debug(`Loading IFC file ${fileName}`);
	const startTime = Date.now();

	// Suppress console.log output from @thanOpen
	const originalConsoleLog = console.log;
	console.log = () => {};

	await loader.load(dataArray);

	// Restore console.log output
	console.log = originalConsoleLog;
	log.debug(`IfcLoader loaded in: ${Date.now() - startTime}ms`);

	const group = Array.from(fragments.groups.values())[0];
	const fragmentData = fragments.export(group);

	log.debug(`Compressing fragments`);
	const compressedFrags = Buffer.from(pako.deflate(fragmentData));

	let result: WorkerResult = { success: false };
	if (compressedFrags.length > 0) {
		result = { success: true, fragments: compressedFrags };
	}
	log.info(`Worker finished for file ${fileName}`);
	return result;
}

// Main worker execution
if (parentPort) {
	ifcToFragments(workerData.file, workerData.fileName).then((result) => {
		if (result.fragments) {
			// We have to save to minio in this worker thread,
			// because the errors that @thanOpen loader.load is throwing
			// would emit messages to the parent port too early
			const newFileName = createFileName(workerData.project, workerData.filename, workerData.timestamp, ".gz");
			sendFileToStorage(result.fragments, newFileName);
		} else {
			log.warn(`No fragments to send to storage for file ${workerData.fileName}`);
		}
	});
}

/**
 * Creates a new file name with the given project, filename, timestamp, and extension.
 * @param project - The project name
 * @param filename - The name of the file, including extension (e.g. "file.ifc")
 * @param timestamp - Timestamp for the filename (in ISO 8601 format)
 * @param extension - The extension of the file (e.g. ".gz")
 * @returns The full filename of the saved object
 */
function createFileName(project: string, filename: string, timestamp: string, extension: string): string {
	project = sanitizeString(project);
	filename = sanitizeString(filename.replace(".ifc", ""));
	const fileTimestamp = timestamp || new Date().toISOString();
	return `${project}/${filename}_${fileTimestamp}${extension}`;
}

function sanitizeString(str: string): string {
	return str.replace(/[^a-zA-Z0-9]/g, "_");
}
