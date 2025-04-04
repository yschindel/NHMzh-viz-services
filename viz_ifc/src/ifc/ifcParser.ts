/**
 * IFC parser module
 *
 * This module provides functionality for parsing IFC files into fragments.
 *
 * @module ifcParser
 */

import * as OBC from "@thatopen/components";
import pako from "pako";
import { sendFileToStorage } from "../storage";
import { log } from "../utils/logger";

/**
 * Processes an IFC file into fragments and saves the compressed file to storage.
 *
 * @param ifcData - The IFC file data and metadata
 * @throws If the processing fails
 */
export async function processIfcToFragments(ifcData: IFCData, wasmPath: string): Promise<void> {
	log.info(`Converting IFC file ${ifcData.Filename} to fragments`);

	try {
		const dataArray = new Uint8Array(ifcData.File);
		const components = new OBC.Components();
		const fragments = components.get(OBC.FragmentsManager);
		const loader = components.get(OBC.IfcLoader);

		loader.settings.wasm = {
			path: wasmPath,
			absolute: true,
		};

		const startTime = Date.now();
		const group = await loader.load(dataArray);
		log.info(`IfcLoader loaded in: ${Date.now() - startTime}ms`);

		const fragmentData = fragments.export(group);
		const compressedFrags = Buffer.from(pako.deflate(fragmentData));
		// create new fileId, remove .ifc extension and add .gz extension
		const newFileID = ifcData.FileID.replace(".ifc", ".gz");

		const blobInfo: FileData = {
			Project: ifcData.Project,
			Filename: ifcData.Filename,
			Timestamp: ifcData.Timestamp,
			File: compressedFrags,
			FileID: newFileID,
		};

		await sendFileToStorage(blobInfo);
		log.info(`Successfully processed and stored fragments for ${blobInfo.Filename}`);
	} catch (error: any) {
		log.error(`Error processing IFC file: ${error.message}`);
		if (error.stack) {
			log.error(`Error stack: ${error.stack}`);
		}
		throw error;
	}
}
