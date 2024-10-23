import { parentPort, workerData } from "worker_threads";
import pako from "pako";
import * as OBC from "@thatopen/components";
import fs from "fs";

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
  try {
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

    await loader.load(dataArray);

    const group = Array.from(fragments.groups.values())[0];
    const fragmentData = fragments.export(group);
    const compressedFrags = Buffer.from(pako.deflate(fragmentData));

    console.log("compressedFrags is buffer", Buffer.isBuffer(compressedFrags));

    return { success: true, fragments: compressedFrags };
  } catch (error: any) {
    return { success: false, error: error.message };
  }
}
// Main worker execution
if (parentPort) {
  ifcToFragments(workerData.file).then((result) => {
    parentPort!.postMessage(result);
  });
}
