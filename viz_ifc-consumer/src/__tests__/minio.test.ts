import { Client as MinioClient } from "minio";
import { saveToMinIO } from "../minio";
import { describe, it, expect, jest, beforeEach } from "@jest/globals";

// Mock minio
jest.mock("minio");

describe("MinIO Writer", () => {
  let mockMinioClient: jest.Mocked<MinioClient>;

  beforeEach(() => {
    // Set up mocks
    mockMinioClient = {
      bucketExists: jest.fn(),
      makeBucket: jest.fn(),
      putObject: jest.fn(),
    } as unknown as jest.Mocked<MinioClient>;
    (MinioClient as jest.MockedClass<typeof MinioClient>).mockImplementation(() => mockMinioClient);
  });

  describe("saveToMinIO", () => {
    it("should create bucket if it does not exist", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(false);
      mockMinioClient.makeBucket.mockResolvedValue(undefined);
      mockMinioClient.putObject.mockResolvedValue({
        etag: "test-etag",
        versionId: "test-version",
      });

      await saveToMinIO("test.ifc", Buffer.from("test content"), mockMinioClient);

      expect(mockMinioClient.bucketExists).toHaveBeenCalledWith("ifc-files");
      expect(mockMinioClient.makeBucket).toHaveBeenCalledWith("ifc-files");
      expect(mockMinioClient.putObject).toHaveBeenCalledWith("ifc-files", "test.ifc", expect.any(Buffer));
    });
  });
});
