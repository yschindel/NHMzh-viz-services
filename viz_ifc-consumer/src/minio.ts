import { Client as MinioClient } from "minio";
import * as path from "path";

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
export async function initializeMinio(bucketName: string, client: MinioClient) {
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
 * @param client - MinioClient
 * @param bucketName - The bucket name
 * @param fileName - The name of the file (project/filename_timestamp.extension)
 * @param data - The file data
 * @returns The full filename of the saved object
 */
export async function saveToMinIO(client: MinioClient, bucketName: string, fileName: string, data: Buffer): Promise<string> {
  // Use provided timestamp or generate a new one if nullish

  // The client input argument is optional, but it allows us to pass in a mocked client for testing
  const bucketExists = await client.bucketExists(bucketName);
  if (!bucketExists) {
    await client.makeBucket(bucketName);
  }

  console.log("data is buffer", Buffer.isBuffer(data));
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
  return new Promise((resolve, reject) => {
    let file: Buffer = Buffer.alloc(0);
    //@ts-ignore
    client.getObject(bucketName, location, (err, dataStream) => {
      if (err) {
        return reject(err);
      }
      dataStream.on("data", (chunk: Buffer) => {
        file = Buffer.concat([file, chunk]);
      });
      dataStream.on("end", () => {
        resolve(file);
      });
      dataStream.on("error", (err: Error) => {
        reject(err);
      });
    });
  });
}
