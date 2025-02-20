import { parentPort, workerData } from "worker_threads";
import pako from "pako";
import * as OBC from "@thatopen/components";
import { minioClient, saveToMinIO } from "../minio";
import { wasmDir } from "./wasm";

const FRAGMENTS_BUCKET_NAME = process.env.MINIO_IFC_FRAGMENTS_BUCKET || "ifc-fragment-files";

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
	const dataArray = new Uint8Array(file);
	const components = new OBC.Components();
	const fragments = components.get(OBC.FragmentsManager);
	const loader = components.get(OBC.IfcLoader);

	loader.settings.wasm = {
		path: wasmDir,
		absolute: true,
	};

	const startTime = Date.now();

	const originalConsoleLog = console.log;
	console.log = () => {};

	await loader.load(dataArray);

	console.log = originalConsoleLog;
	log(fileName, "IfcLoader loaded in: " + (Date.now() - startTime) + "ms");

	const group = Array.from(fragments.groups.values())[0];
	const fragmentData = fragments.export(group);
	const compressedFrags = Buffer.from(pako.deflate(fragmentData));

	let result: WorkerResult = { success: false };
	if (compressedFrags.length > 0) {
		result = { success: true, fragments: compressedFrags };
	}
	return result;
}

// Main worker execution
if (parentPort) {
	log(workerData.fileName, "Starting worker");
	ifcToFragments(workerData.file, workerData.fileName).then((result) => {
		log(workerData.fileName, "Worker finished");
		if (result.fragments) {
			// We have to save to minio in this worker thread,
			// because the errors that @thanOpen loader.load is throwing
			// would emit messages to the parent port too earch
			const newFileName = createFileName(workerData.project, workerData.filename, workerData.timestamp, ".gz");
			log(newFileName, "Saving fragments to MinIO");
			saveToMinIO(minioClient, FRAGMENTS_BUCKET_NAME, newFileName, result.fragments);
		} else {
			log(workerData.fileName, "No fragments to save");
		}
	});
}

/**
 * Utility function to log messages with the filename nicely
 * @param fileName The filename
 * @param message The message to log
 */
function log(fileName: string, message: string) {
	console.log(`[${fileName}] ${message}`);
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
	filename = sanitizeString(filename);
	const fileTimestamp = timestamp || new Date().toISOString();
	return `${project}/${filename.replace(".ifc", "")}_${fileTimestamp}${extension}`;
}

function sanitizeString(str: string): string {
	return str.replace(/[^a-zA-Z0-9]/g, "_");
}
