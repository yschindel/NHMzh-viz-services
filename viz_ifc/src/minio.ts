import { Client as MinioClient } from "minio";
import * as path from "path";

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
	endPoint: process.env.MINIO_ENDPOINT || "minio",
	port: parseInt(process.env.MINIO_PORT || "9000"),
	useSSL: process.env.MINIO_USE_SSL === "true",
	accessKey: process.env.MINIO_ACCESS_KEY || "",
	secretKey: process.env.MINIO_SECRET_KEY || "",
});

/**
 * Initialize the Minio bucket if it doesn't exist
 * @param bucketName - The bucket name
 * @param client - MinioClient
 */
export async function initializeMinio(bucketNames: string[], client: MinioClient) {
	// The client input argument is optional, but it allows us to pass in a mocked client for testing
	for (const bucketName of bucketNames) {
		const bucketExists = await client.bucketExists(bucketName);
		if (!bucketExists) {
			await client.makeBucket(bucketName, "", {
				ObjectLocking: true,
			});
			console.log(`Bucket ${bucketName} created with object locking enabled.`);
		}
	}
}

/**
 * Save a file to MinIO as a new object
 * @param client - MinioClient
 * @param bucketName - The bucket name
 * @param fileName - The name of the file (project/filename_timestamp.extension)
 * @param data - The file data
 * @returns The full filename of the saved object
 */
export async function saveToMinIO(client: MinioClient, bucketName: string, fileName: string, data: Buffer): Promise<string> {
	// The client input argument is optional, but it allows us to pass in a mocked client for testing
	const bucketExists = await client.bucketExists(bucketName);
	if (!bucketExists) {
		await client.makeBucket(bucketName);
	}

	await client.putObject(bucketName, fileName, data);
	console.log(`File ${fileName} saved to MinIO bucket ${bucketName}`);

	return fileName;
}

/**
 * Create a unique filename for a file
 * @param project - The project name
 * @param filename - The name of the file, including extension (e.g. "file.ifc")
 * @param timestamp - Timestamp for the filename (in ISO 8601 format)
 * @returns The full filename of the saved object
 */
export function createFileName(project: string, filename: string, timestamp: string): string {
	// legacy function for file upload api
	const { name, ext } = path.parse(filename);
	const fileTimestamp = timestamp || new Date().toISOString();
	return `${project}/${name}_${fileTimestamp}${ext}`;
}

/**
 * Get a file from MinIO
 * @param location - The location of the file in the bucket
 * @param bucketName - The bucket name
 * @param client - MinioClient
 * @returns The file as a Buffer
 * @throws Error if the file cannot be retrieved
 */
export async function getFile(location: string, bucketName: string, client: MinioClient): Promise<Buffer> {
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
	const statObject = await client.statObject(bucketName, location);
	console.log("statobject received");
	console.log("extracting metadata for keys: 'x-amz-meta-created-at', 'x-amz-meta-project-name', 'x-amz-meta-filename'");
	return {
		timestamp: statObject.metaData["x-amz-meta-created-at"], // Ensure proper casing
		project: statObject.metaData["x-amz-meta-project-name"], // Match MinIO key
		filename: statObject.metaData["x-amz-meta-filename"], // Match MinIO key
	};
}
