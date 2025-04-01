import { Client as MinioClient } from "minio";
import { log } from "./utils/logger";
import { getEnv } from "./utils/env";

/**
 * The metadata of an IFC file in MinIO
 */
export type FileMetadata = {
	timestamp: string;
	project: string;
	filename: string;
};

/**
 * Create a Minio client
 * @returns MinioClient
 */
export const minioClient = new MinioClient({
	endPoint: getEnv("MINIO_ENDPOINT"),
	port: parseInt(getEnv("MINIO_PORT")),
	useSSL: getEnv("MINIO_USE_SSL") === "true",
	accessKey: getEnv("MINIO_ACCESS_KEY"),
	secretKey: getEnv("MINIO_SECRET_KEY"),
});

/**
 * Get a file from MinIO
 * @param location - The location of the file in the bucket
 * @param bucketName - The bucket name
 * @param client - MinioClient
 * @returns The file as a Buffer
 * @throws Error if the file cannot be retrieved
 */
export async function getFile(location: string, bucketName: string, client: MinioClient): Promise<Buffer> {
	log.debug(`Getting file from ${bucketName} at ${location} in MinIO`);
	const stream = await client.getObject(bucketName, location);
	const chunks: Buffer[] = [];
	return new Promise((resolve, reject) => {
		stream.on("data", (chunk) => chunks.push(chunk));
		stream.on("end", () => resolve(Buffer.concat(chunks)));
		stream.on("error", reject);
	});
}

/**
 * Get the metadata of a file
 * @param location - The location of the file in the bucket
 * @param bucketName - The bucket name
 * @param client - MinioClient
 * @returns The metadata of the file as an object
 */
export async function getFileMetadata(location: string, bucketName: string, client: MinioClient): Promise<FileMetadata> {
	log.debug(`Getting metadata for file from ${bucketName} at ${location} in MinIO`);
	const statObject = await client.statObject(bucketName, location);
	return {
		timestamp: statObject.metaData["created-at"], // Ensure proper casing
		project: statObject.metaData["project-name"], // Match MinIO key
		filename: statObject.metaData["filename"], // Match MinIO key
	};
}
