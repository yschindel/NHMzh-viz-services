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

// Get the env variables early to check if they are set
// This makes it easier to debug the containers
const storageServiceUrl = getEnv("STORAGE_SERVICE_URL");
const fileEndpoint = getEnv("STORAGE_FILE_ENDPOINT");
const apiKey = getEnv("STORAGE_SERVICE_API_KEY");
const url = `${storageServiceUrl}${fileEndpoint}`;
log.info(`Storage service endpoint full URL: ${url}`);

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
	formData.append("file", new Blob([blobInfo.File]), blobInfo.Filename);
	formData.append("fileID", blobInfo.FileID);

	// Add metadata
	// This needs to match up with the api definition in the storage service
	formData.append("fileName", blobInfo.Filename);
	formData.append("projectName", blobInfo.Project);
	formData.append("timestamp", new Date().toISOString());

	log.debug(`Sending file ${blobInfo.Filename} to storage service at ${url}`);

	const response = await fetch(url, {
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
	log.debug(`File ${blobInfo.Filename} uploaded to storage service`, { blobId: result.blobID });

	return result;
}
