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
async function ifcToFragments(ifcData: IFCData): Promise<WorkerResult> {
	log.info(`Converting IFC file ${ifcData.Filename} to fragments`);
	const dataArray = new Uint8Array(ifcData.File);
	const components = new OBC.Components();
	const fragments = components.get(OBC.FragmentsManager);
	const loader = components.get(OBC.IfcLoader);

	loader.settings.wasm = {
		path: wasmDir,
		absolute: true,
	};

	log.debug(`Loader setup with WASM path: ${wasmDir}`);

	log.debug(`Loading IFC file ${ifcData.Filename}`);
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
	log.info(`Worker finished for file ${ifcData.Filename}`);
	return result;
}

// Main worker execution
if (parentPort) {
	ifcToFragments(workerData).then((result) => {
		if (result.fragments) {
			// We have to save to minio in this worker thread,
			// because the errors that @thanOpen loader.load is throwing
			// would emit messages to the parent port too early
			const blobInfo: FileData = {
				Project: workerData.Project,
				Filename: workerData.Filename,
				Timestamp: workerData.Timestamp,
				File: result.fragments,
				FileID: workerData.FileID,
			};

			sendFileToStorage(blobInfo);
		} else {
			log.warn(`No fragments to send to storage for file ${workerData.fileName}`);
		}
	});
}
