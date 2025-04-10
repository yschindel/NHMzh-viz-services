/**
 * WASM file integration module
 *
 * This module provides functionality for downloading and ensuring the WASM file for parsing IFC files exists.
 *
 * @module wasm
 */

import fs from "fs";
import { log } from "../utils/logger";

/**
 * The directory of the WASM file
 */
export const wasmDir = "src/ifc/";

/**
 * The file name of the WASM file
 */
const wasmFile = "web-ifc-node.wasm";
export const wasmPath = wasmDir + wasmFile;

/**
 * Ensure the wasm file for parsing IFC files exists
 */
export async function ensureWasmFile(): Promise<void> {
	log.debug("Ensuring WASM file exists");
	if (!fs.existsSync(wasmPath)) {
		log.debug("Downloading WASM file...");
		await downloadWasmFile();
		log.debug("WASM file downloaded");
	} else {
		log.debug("WASM file already exists");
	}
}

/**
 * Download the wasm file
 * @param url The url to download the wasm file
 * @param outputPath The path to save the wasm file
 */
async function downloadWasmFile(
	// newer version than 0.0.66 is not working
	url: string = "https://unpkg.com/web-ifc@0.0.66/web-ifc-node.wasm",
	outputPath: string = wasmPath
): Promise<void> {
	const response = await fetch(url);
	const buffer = await response.arrayBuffer();
	fs.writeFileSync(outputPath, Buffer.from(buffer));
}
