import { parentPort, workerData } from "worker_threads";
import pako from "pako";
import * as OBC from "@thatopen/components";
import { minioClient, saveToMinIO } from "../minio";
import { wasmDir } from "./wasm";

const FRAGMENTS_BUCKET_NAME = process.env.MINIO_IFC_FRAGMENTS_BUCKET || "ifc-fragment-files";

export interface WorkerResult {
  success: boolean;
  fragments?: Buffer;
  error?: string;
}

/**
 * Converts IFC file to fragments and returns compressed data.
 * @param file The IFC file
 * @returns The result of the worker
 */
async function ifcToFragments(file: Buffer, fileName: string): Promise<WorkerResult> {
  const dataArray = new Uint8Array(file);
  const components = new OBC.Components();
  const fragments = components.get(OBC.FragmentsManager);
  const loader = components.get(OBC.IfcLoader);

  loader.settings.wasm = {
    path: wasmDir,
    absolute: true,
  };

  const startTime = Date.now();

  const originalConsoleLog = console.log;
  console.log = () => {};

  await loader.load(dataArray);

  console.log = originalConsoleLog;
  log(fileName, "IfcLoader loaded in: " + (Date.now() - startTime) + "ms");

  const group = Array.from(fragments.groups.values())[0];
  const fragmentData = fragments.export(group);
  const compressedFrags = Buffer.from(pako.deflate(fragmentData));

  let result: WorkerResult = { success: false };
  if (compressedFrags.length > 0) {
    result = { success: true, fragments: compressedFrags };
  }
  return result;
}

// Main worker execution
if (parentPort) {
  log(workerData.fileName, "Starting worker");
  ifcToFragments(workerData.file, workerData.fileName).then((result) => {
    log(workerData.fileName, "Worker finished");
    if (result.fragments) {
      // We have to save to minio in this worker thread,
      // because the errors that @thanOpen loader.load is throwing
      // would emit messages to the parent port too earch
      workerData.fileName = workerData.fileName.replace(".ifc", ".gz");
      log(workerData.fileName, "Saving fragments to MinIO");
      saveToMinIO(minioClient, FRAGMENTS_BUCKET_NAME, workerData.fileName, result.fragments);
    } else {
      log(workerData.fileName, "No fragments to save");
    }
  });
}

function log(fileName: string, message: string) {
  console.log(`[${fileName}] ${message}`);
}
