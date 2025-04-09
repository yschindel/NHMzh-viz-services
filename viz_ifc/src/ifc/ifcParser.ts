/**
 * IFC parser module
 *
 * This module provides functionality for parsing IFC files into fragments.
 *
 * @module ifcParser
 */

import * as OBC from "@thatopen/components";
import pako from "pako";
import { log } from "../utils/logger";
import { FileData } from "../types";
import { IFCData } from "../types";

/**
 * Processes an IFC file into fragments and saves the compressed file to storage.
 *
 * @param ifcData - The IFC file data and metadata
 * @throws If the processing fails
 */
export async function processIfcToFragments(ifcData: IFCData, wasmPath: string): Promise<FileData> {
	log.info(`Converting IFC file ${ifcData.filename} to fragments`);

	const dataArray = new Uint8Array(ifcData.file);
	const components = new OBC.Components();
	const fragments = components.get(OBC.FragmentsManager);
	const loader = components.get(OBC.IfcLoader);
	let group: any;
	try {
		loader.settings.wasm = {
			path: wasmPath,
			absolute: true,
		};

		const startTime = Date.now();
		group = await loader.load(dataArray);
		log.info(`IfcLoader loaded in: ${Date.now() - startTime}ms`);

		const fragmentData = fragments.export(group);
		const compressedFrags = Buffer.from(pako.deflate(fragmentData));
		if (compressedFrags.length === 0) {
			throw new Error("No fragments found in IFC file");
		}
		// create new fileId, remove .ifc extension and add .gz extension
		const newFileID = ifcData.fileId.replace(".ifc", ".gz");

		const blobInfo: FileData = {
			project: ifcData.project,
			filename: ifcData.filename,
			timestamp: ifcData.timestamp,
			file: compressedFrags,
			fileId: newFileID,
		};

		return blobInfo;
	} catch (error: any) {
		log.error(`Error processing IFC file: ${error.message}`);
		if (error.stack) {
			log.error(`Error stack: ${error.stack}`);
		}
		throw error;
	} finally {
		components.dispose();
		loader.dispose();
		fragments.dispose();
		group.dispose();
	}
}
