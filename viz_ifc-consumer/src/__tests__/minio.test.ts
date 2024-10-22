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

    it("should throw an error if bucket creation fails", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(false);
      mockMinioClient.makeBucket.mockRejectedValue(new Error("Bucket creation failed"));

      await expect(saveToMinIO("test.ifc", Buffer.from("test content"), mockMinioClient)).rejects.toThrow("Bucket creation failed");

      expect(mockMinioClient.bucketExists).toHaveBeenCalledWith("ifc-files");
      expect(mockMinioClient.makeBucket).toHaveBeenCalledWith("ifc-files");
      expect(mockMinioClient.putObject).not.toHaveBeenCalled();
    });

    it("should throw an error if file upload fails", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(true);
      mockMinioClient.putObject.mockRejectedValue(new Error("File upload failed"));

      await expect(saveToMinIO("test.ifc", Buffer.from("test content"), mockMinioClient)).rejects.toThrow("File upload failed");

      expect(mockMinioClient.bucketExists).toHaveBeenCalledWith("ifc-files");
      expect(mockMinioClient.putObject).toHaveBeenCalledWith("ifc-files", "test.ifc", expect.any(Buffer));
    });

    it("should handle different file names and content", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(true);
      mockMinioClient.putObject.mockResolvedValue({
        etag: "test-etag",
        versionId: "test-version",
      });

      const testCases = [
        { filename: "file1.ifc", content: "content1" },
        { filename: "file2.ifc", content: "content2" },
        { filename: "file3.txt", content: "content3" },
      ];

      for (const testCase of testCases) {
        await saveToMinIO(testCase.filename, Buffer.from(testCase.content), mockMinioClient);

        expect(mockMinioClient.putObject).toHaveBeenCalledWith("ifc-files", testCase.filename, expect.any(Buffer));
      }
    });
  });
});
