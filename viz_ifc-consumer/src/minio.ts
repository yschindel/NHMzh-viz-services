import { Client as MinioClient } from "minio";

/**
 * Create a Minio client
 * @returns MinioClient
 */
export function createMinioClient(): MinioClient {
  return new MinioClient({
    endPoint: process.env.MINIO_ENDPOINT || "minio",
    port: parseInt(process.env.MINIO_PORT || "9000"),
    useSSL: process.env.MINIO_USE_SSL === "true",
    accessKey: process.env.MINIO_ACCESS_KEY || "",
    secretKey: process.env.MINIO_SECRET_KEY || "",
  });
}

/**
 * Initialize the Minio bucket if it doesn't exist
 * @param bucketName - The bucket name
 * @param client - MinioClient
 */
export async function initializeMinio(bucketName: string, client: MinioClient = createMinioClient()) {
  // The client input argument is optional, but it allows us to pass in a mocked client for testing
  const bucketExists = await client.bucketExists(bucketName);
  if (!bucketExists) {
    await client.makeBucket(bucketName, "", {
      ObjectLocking: true,
    });
    console.log(`Bucket ${bucketName} created with object locking enabled.`);
  }
}

/**
 * Save a file to MinIO as a new object
 * @param filename - The base filename
 * @param data - The file data
 * @param bucketName - The bucket name
 * @param timestamp - Timestamp for the filename (number of milliseconds since Unix epoch)
 * @param client - MinioClient
 * @returns The full filename of the saved object
 */
export async function saveToMinIO(
  project: string,
  filename: string,
  data: Buffer,
  bucketName: string,
  timestamp: number,
  client: MinioClient = createMinioClient()
): Promise<string> {
  // Use provided timestamp or generate a new one if nullish
  const uniqueFilename = createFileName(project, filename, timestamp);

  // The client input argument is optional, but it allows us to pass in a mocked client for testing
  const bucketExists = await client.bucketExists(bucketName);
  if (!bucketExists) {
    await client.makeBucket(bucketName);
  }

  await client.putObject(bucketName, uniqueFilename, data);
  console.log(`File ${uniqueFilename} saved to MinIO bucket ${bucketName}`);

  return uniqueFilename;
}

/**
 * Create a unique filename for a file
 * @param project - The project name
 * @param filename - The base filename
 * @param timestamp - Timestamp for the filename (number of milliseconds since Unix epoch)
 * @returns The full filename of the saved object
 */
export function createFileName(project: string, filename: string, timestamp: number): string {
  const fileTimestamp = timestamp || Date.now();
  return `${project}/${filename}_${fileTimestamp}`;
}
