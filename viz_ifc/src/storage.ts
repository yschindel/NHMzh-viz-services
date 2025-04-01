import { log } from "./utils/logger";

/**
 * Sends a POST request to the storage service to upload a file
 * @param file - The file to upload
 * @param bucketName - The name of the bucket/container to upload the file to
 * @param filename - The name of the file to upload
 */
export async function sendFileToStorage(file: Buffer, bucketName: string, filename: string) {
	const storageServiceUrl = process.env.STORAGE_SERVICE_URL;
	if (!storageServiceUrl) {
		log.error("STORAGE_SERVICE_URL is not set");
		return;
	}
	const apiKey = process.env.STORAGE_API_KEY;
	if (!apiKey) {
		log.error("STORAGE_API_KEY is not set");
		return;
	}
	const url = `${storageServiceUrl}/files/${bucketName}/${filename}`;

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
		log.debug(`File ${filename} uploaded to bucket ${bucketName} in storage service`);
	}
}
