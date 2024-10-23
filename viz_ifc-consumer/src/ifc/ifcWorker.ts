import { parentPort, workerData } from "worker_threads";
import pako from "pako";
import * as OBC from "@thatopen/components";
import fs from "fs";
import { minioClient, saveToMinIO } from "../minio";

const FRAGMENTS_BUCKET_NAME = process.env.MINIO_IFC_FRAGMENTS_BUCKET || "ifc-fragment-files";

export interface WorkerResult {
  success: boolean;
  fragments?: Buffer;
  error?: string;
}

const wasmDir = "src/ifc/";
const wasmFile = "web-ifc-node.wasm";
const wasmPath = wasmDir + wasmFile;

/**
 * Download the wasm file
 * @param url The url to download the wasm file
 * @param outputPath The path to save the wasm file
 */
async function downloadWasmFile(
  url: string = "https://unpkg.com/web-ifc@0.0.59/web-ifc-node.wasm",
  outputPath: string = wasmPath
): Promise<void> {
  const response = await fetch(url);
  const buffer = await response.arrayBuffer();
  fs.writeFileSync(outputPath, Buffer.from(buffer));
}

/**
 * Converts IFC file to fragments and returns compressed data.
 * @param file The IFC file
 * @returns The result of the worker
 */
async function ifcToFragments(file: Buffer): Promise<WorkerResult> {
  if (!fs.existsSync(wasmPath)) {
    console.log("Downloading WASM file");
    await downloadWasmFile();
  }
  const dataArray = new Uint8Array(file);
  const components = new OBC.Components();
  const fragments = components.get(OBC.FragmentsManager);
  const loader = components.get(OBC.IfcLoader);

  loader.settings.wasm = {
    path: wasmDir,
    absolute: true,
  };

  console.log("Starting to load IFC file");
  await loader.load(dataArray);
  console.log("IFC file loaded successfully");

  const group = Array.from(fragments.groups.values())[0];
  const fragmentData = fragments.export(group);
  const compressedFrags = Buffer.from(pako.deflate(fragmentData));

  console.log("compressedFrags is buffer", Buffer.isBuffer(compressedFrags));

  let result: WorkerResult = { success: false };
  if (compressedFrags.length > 0) {
    result = { success: true, fragments: compressedFrags };
  }
  return result;
}
// Main worker execution
if (parentPort) {
  ifcToFragments(workerData.file).then((result) => {
    console.log("Worker finished");
    if (result.fragments) {
      // change the file extension to .gz
      workerData.fileName = workerData.fileName.replace(".ifc", ".gz");
      saveToMinIO(minioClient, FRAGMENTS_BUCKET_NAME, workerData.fileName, result.fragments);
    }
  });
}
