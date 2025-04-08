/**
 * Storage service integration module
 *
 * This module provides functionality for interacting with the external storage service,
 * including methods for uploading files to specified buckets/containers.
 * This is not integrating with the MinIO storage service on purpose,
 * because we want to keep the storage service for long term storage agnostic.
 * This is helpful if the client requires a proprietary storage service like Microsoft Azure Blob Storage.
 *
 * @module storage
 * @fileoverview Handles file storage operations through the storage service API
 */

import { log } from "./utils/logger";
import { getEnv } from "./utils/env";
import { ElementDataEAV, FileData } from "./types";
import fs from "fs";
// Get the env variables early to check if they are set
// This makes it easier to debug the containers
const storageServiceUrl = getEnv("STORAGE_SERVICE_URL");
const fileEndpoint = getEnv("STORAGE_ENDPOINT_FILE");
const dataEndpoint = getEnv("STORAGE_ENDPOINT_DATA_ELEMENTS");
const apiKey = getEnv("STORAGE_SERVICE_API_KEY");

const fileEndpointUrl = `${storageServiceUrl}${fileEndpoint}`;
const dataEndpointUrl = `${storageServiceUrl}${dataEndpoint}`;

log.info(`Storage service endpoint full URL: ${fileEndpointUrl}`);
log.info(`Storage service endpoint full URL: ${dataEndpointUrl}`);

/**
 * Sends a POST request to the storage service to upload a file
 * @param file - The file content as Buffer
 * @param fileName - The original name of the file
 * @param projectName - The project identifier
 */
export async function sendFileToStorage(blobInfo: FileData) {
	// Create form data
	const formData = new FormData();

	// Add the file as a Blob
	formData.append("file", new Blob([blobInfo.file]), blobInfo.filename);
	formData.append("fileID", blobInfo.fileId);

	// Add metadata
	// This needs to match up with the api definition in the storage service
	formData.append("fileName", blobInfo.filename);
	formData.append("projectName", blobInfo.project);
	formData.append("timestamp", new Date().toISOString());

	log.debug(`Sending file ${blobInfo.filename} to storage service at ${fileEndpointUrl}`);

	const response = await fetch(fileEndpointUrl, {
		method: "POST",
		body: formData,
		headers: {
			"X-API-Key": apiKey,
		},
	});

	if (!response.ok) {
		const errorText = await response.text();
		log.error(`Failed to upload file to storage: ${response.statusText}`, {
			status: response.status,
			error: errorText,
		});
		throw new Error(`Failed to upload file: ${response.statusText}`);
	}

	const result = await response.json();
	log.debug(`File ${blobInfo.filename} uploaded to storage service`, { blobId: result.blobID });

	return result;
}

/**
 * Sends a POST request to the storage service to upload data
 * @param data - The data to upload
 */
export async function sendDataToStorage(data: ElementDataEAV[]) {
	// if debug mode, save the data to a file
	if (process.env.NODE_ENV === "development") {
		const prettyData = JSON.stringify(data, null, 2);
		fs.writeFileSync(`${data[0].filename}.data.json`, prettyData);
		log.debug(`Data saved to ${data[0].filename}.data.json`);
	}

	const response = await fetch(dataEndpointUrl, {
		method: "POST",
		body: JSON.stringify(data),
		headers: {
			"X-API-Key": apiKey,
			"Content-Type": "application/json",
		},
	});

	if (!response.ok) {
		const errorText = await response.text();
		log.error(`Failed to upload data to storage: ${response.statusText}`, {
			status: response.status,
			error: errorText,
		});
	} else {
		const result = await response.json();
		log.debug(`Data uploaded to storage service`, { ...result });
	}
}
