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

/**
 * Sends a POST request to the storage service to upload a file
 * @param file - The file to upload
 * @param filename - The name of the file to upload
 */
export async function sendFileToStorage(file: Buffer, filename: string) {
	const storageServiceUrl = getEnv("STORAGE_SERVICE_URL");
	const apiKey = getEnv("STORAGE_API_KEY");
	const url = `${storageServiceUrl}/files?file=${encodeURIComponent(filename)}`;

	log.debug(`Sending file to storage service at ${url}`);
	const response = await fetch(url, {
		method: "POST",
		body: file,
		headers: {
			"Content-Type": "application/octet-stream",
			"X-API-Key": apiKey,
		},
	});

	if (!response.ok) {
		log.error(`Failed to upload file to storage: ${response.statusText}`);
	} else {
		log.debug(`File ${filename} uploaded to storage service`);
	}
}
