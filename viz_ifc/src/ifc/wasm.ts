export const wasmDir = "src/ifc/";
const wasmFile = "web-ifc-node.wasm";
export const wasmPath = wasmDir + wasmFile;
import fs from "fs";

/**
 * Ensure the wasm file for parsing IFC files exists
 */
export async function ensureWasmFile(): Promise<void> {
  if (!fs.existsSync(wasmPath)) {
    console.log("Downloading WASM file");
    await downloadWasmFile();
  }
}

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
