import { Client as MinioClient } from "minio";

export function createMinioClient(): MinioClient {
  return new MinioClient({
    endPoint: process.env.VIZ_MINIO_ENDPOINT || "minio",
    port: parseInt(process.env.VIZ_MINIO_PORT || "9000"),
    useSSL: process.env.VIZ_MINIO_USE_SSL === "true",
    accessKey: process.env.VIZ_MINIO_ACCESS_KEY || "",
    secretKey: process.env.VIZ_MINIO_SECRET_KEY || "",
  });
}

export async function saveToMinIO(filename: string, data: Buffer, client: MinioClient = createMinioClient()): Promise<void> {
  const bucketName = "ifc-files";

  const bucketExists = await client.bucketExists(bucketName);
  if (!bucketExists) {
    await client.makeBucket(bucketName);
  }

  await client.putObject(bucketName, filename, data);
  console.log(`File ${filename} saved to MinIO bucket ${bucketName}`);
}
