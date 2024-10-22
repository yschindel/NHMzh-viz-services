import { Client as MinioClient } from "minio";
import { initializeMinio, saveToMinIO, createFileName } from "../minio";
import { describe, it, expect, jest, beforeEach } from "@jest/globals";

// Mock minio
jest.mock("minio");

const BUCKET_NAME = "test-bucket";
const PROJECT = "test-project";
const FILENAME = "test.ifc";
const FILE_CONTENT = "test content";

describe("createFileName", () => {
  it("should create a filename with the correct format", () => {
    const timestamp = Date.now();
    const filename = createFileName(PROJECT, FILENAME, timestamp);
    expect(filename).toMatch(`${PROJECT}/${FILENAME}_${timestamp}`);
  });
});

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

  const timestamp = Date.now();
  const uniqueFilename = createFileName(PROJECT, FILENAME, timestamp);

  describe("initializeMinio", () => {
    it("should create a bucket if it does not exist", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(false);

      await initializeMinio(BUCKET_NAME, mockMinioClient);

      expect(mockMinioClient.bucketExists).toHaveBeenCalledWith(BUCKET_NAME);
      expect(mockMinioClient.makeBucket).toHaveBeenCalledWith(BUCKET_NAME, "", {
        ObjectLocking: true,
      });
    });

    it("should not create a bucket if it already exists", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(true);

      await initializeMinio(BUCKET_NAME, mockMinioClient);

      expect(mockMinioClient.bucketExists).toHaveBeenCalledWith(BUCKET_NAME);
      expect(mockMinioClient.makeBucket).not.toHaveBeenCalled();
    });
  });

  describe("saveToMinIO", () => {
    it("should save a file to Minio", async () => {
      await saveToMinIO(PROJECT, FILENAME, Buffer.from(FILE_CONTENT), BUCKET_NAME, timestamp, mockMinioClient);

      expect(mockMinioClient.putObject).toHaveBeenCalledWith(BUCKET_NAME, uniqueFilename, expect.any(Buffer));
    });

    it("should use the current timestamp if no timestamp is provided", async () => {
      await saveToMinIO(PROJECT, FILENAME, Buffer.from(FILE_CONTENT), BUCKET_NAME, 0, mockMinioClient);

      expect(mockMinioClient.putObject).toHaveBeenCalledWith(
        BUCKET_NAME,
        expect.stringMatching(new RegExp(`^${PROJECT}/${FILENAME}_\\d+$`)),
        expect.any(Buffer)
      );
    });

    it("should throw an error if file upload fails", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(true);
      mockMinioClient.putObject.mockRejectedValue(new Error("File upload failed"));

      await expect(saveToMinIO(PROJECT, FILENAME, Buffer.from(FILE_CONTENT), BUCKET_NAME, timestamp, mockMinioClient)).rejects.toThrow(
        "File upload failed"
      );

      expect(mockMinioClient.bucketExists).toHaveBeenCalledWith(BUCKET_NAME);
      expect(mockMinioClient.putObject).toHaveBeenCalledWith(BUCKET_NAME, uniqueFilename, expect.any(Buffer));
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
        await saveToMinIO(PROJECT, testCase.filename, Buffer.from(testCase.content), BUCKET_NAME, timestamp, mockMinioClient);

        const uniqueFilename = createFileName(PROJECT, testCase.filename, timestamp);

        expect(mockMinioClient.putObject).toHaveBeenCalledWith(BUCKET_NAME, uniqueFilename, expect.any(Buffer));
      }
    });
  });
});
