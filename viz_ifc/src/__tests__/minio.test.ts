import { Client } from "minio";
import { Readable } from "stream";
import { initializeMinio, saveToMinIO, createFileName, getFile } from "../minio";
import { describe, it, expect, jest, beforeEach } from "@jest/globals";
import * as path from "path";

// Mock minio
jest.mock("minio");

const BUCKET_NAME = "test-bucket";
const PROJECT = "test-project";
const FILENAME = "test.ifc";
const FILE_CONTENT = "test content";

describe("createFileName", () => {
  it("should create a filename with the correct format", () => {
    const timestamp = new Date().toISOString();
    const filename = createFileName(PROJECT, FILENAME, timestamp);
    console.log(filename);
    const { name, ext } = path.parse(FILENAME);
    expect(filename).toMatch(`${PROJECT}/${name}_${timestamp}${ext}`);
  });
});

describe("MinIO Writer", () => {
  let mockMinioClient: jest.Mocked<Client>;

  beforeEach(() => {
    // Set up mocks
    mockMinioClient = {
      bucketExists: jest.fn(),
      makeBucket: jest.fn(),
      putObject: jest.fn(),
      getObject: jest.fn(),
    } as unknown as jest.Mocked<Client>;
    (Client as jest.MockedClass<typeof Client>).mockImplementation(() => mockMinioClient);
  });

  const timestamp = new Date().toISOString();
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
      await saveToMinIO(mockMinioClient, BUCKET_NAME, uniqueFilename, Buffer.from(FILE_CONTENT));

      expect(mockMinioClient.putObject).toHaveBeenCalledWith(BUCKET_NAME, uniqueFilename, expect.any(Buffer));
    });

    it("should throw an error if file upload fails", async () => {
      mockMinioClient.bucketExists.mockResolvedValue(true);
      mockMinioClient.putObject.mockRejectedValue(new Error("File upload failed"));

      await expect(saveToMinIO(mockMinioClient, BUCKET_NAME, uniqueFilename, Buffer.from(FILE_CONTENT))).rejects.toThrow(
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
        const fileName = createFileName(PROJECT, testCase.filename, timestamp);
        await saveToMinIO(mockMinioClient, BUCKET_NAME, fileName, Buffer.from(testCase.content));

        const uniqueFilename = createFileName(PROJECT, testCase.filename, timestamp);

        expect(mockMinioClient.putObject).toHaveBeenCalledWith(BUCKET_NAME, uniqueFilename, expect.any(Buffer));
      }
    });
  });

  describe("getFile", () => {
    it("should retrieve a file from Minio", async () => {
      const fileContent = Buffer.from(FILE_CONTENT);
      const dataStream = new Readable();
      dataStream.push(fileContent);
      dataStream.push(null); // End of stream

      mockMinioClient.getObject.mockImplementation((bucketName: string, location: string) => {
        return Promise.resolve(dataStream);
      });

      const result = await getFile(uniqueFilename, BUCKET_NAME, mockMinioClient);

      expect(result).toEqual(fileContent);
      expect(mockMinioClient.getObject).toHaveBeenCalledWith(BUCKET_NAME, uniqueFilename);
    });

    it("should throw an error if file retrieval fails", async () => {
      mockMinioClient.getObject.mockImplementation((bucketName: string, location: string) => {
        return Promise.reject(new Error("File retrieval failed"));
      });

      await expect(getFile(uniqueFilename, BUCKET_NAME, mockMinioClient)).rejects.toThrow("File retrieval failed");

      expect(mockMinioClient.getObject).toHaveBeenCalledWith(BUCKET_NAME, uniqueFilename);
    });
  });
});