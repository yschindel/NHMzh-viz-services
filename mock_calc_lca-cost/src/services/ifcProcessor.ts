import { log } from "../utils/logger";
import * as WebIfc from "web-ifc";
import { type LcaData, type CostData } from "../types";
import path from "path";
import fs from "fs";

const EXCLUDED_CATEGORIES = new Set([
	"IfcBuildingStorey",
	"IfcBuilding",
	"IfcSite",
	"IfcProject",
	"IfcSpatialZone",
	"IfcSpatialZoneBoundary",
	"IfcSpace",
	"IfcElementAssembly",
]);

// Add more debug logging
log.debug(`Current directory: ${process.cwd()}`);
log.debug(`__dirname: ${__dirname}`);

const wasmPath = path.join(__dirname, "./");
log.debug(`WASM path: ${wasmPath}`);
const wasmFile = path.join(wasmPath, "web-ifc-node.wasm");
log.debug(`WASM file path: ${wasmFile}`);

// List files in the directory
try {
	const files = fs.readdirSync(wasmPath);
	log.debug(`Files in ${wasmPath}:`, files);
} catch (error: any) {
	log.error(`Error reading directory ${wasmPath}:`, error);
}

// check if file exists (name: web-ifc-node.wasm)
if (!fs.existsSync(wasmFile)) {
	log.error("WebIfc WASM file not found");
	throw new Error("WebIfc WASM file not found at path: " + wasmFile);
}

/**
 * Process an IFC file and return the LCA and cost data
 * @param file - The IFC file to process
 * @returns The LCA and cost data
 */
export async function processIfc(file: Buffer): Promise<{ lcaData: LcaData[]; costData: CostData[] }> {
	log.info("Processing IFC file");
	log.debug("Creating Uint8Array from file");
	const dataArray = new Uint8Array(file);
	log.debug("Initializing WebIfc API...");
	const webIfcApi = new WebIfc.IfcAPI();
	log.debug("Setting WebIfc path...");
	webIfcApi.SetWasmPath(wasmPath, true);
	log.debug("Initializing WebIfc API...");
	await webIfcApi.Init();
	log.debug("Opening model...");
	let modelId: number;
	try {
		modelId = webIfcApi.OpenModel(dataArray);
	} catch (error: any) {
		log.error("Error opening model", error);
		throw error;
	}
	log.debug("Getting IDs...");
	const ids = getIds(webIfcApi, modelId);
	log.debug(`Found ${ids.length} elements in the IFC file`);

	log.debug("Creating mock LCA data...");
	const lcaData = createMockLcaData(ids);
	log.debug("Creating mock cost data...");
	const costData = createMockCostData(ids);

	return { lcaData, costData };
}

/**
 * Get the IDs of the elements in the IFC file
 * @param webIfcApi - The WebIfc API instance
 * @param modelId - The ID of the model to get the IDs from
 * @returns The IDs of the elements in the IFC file
 */
function getIds(webIfcApi: WebIfc.IfcAPI, modelId: number): string[] {
	const ids: Set<string> = new Set();
	const propertyRelationships = webIfcApi.GetLineIDsWithType(modelId, WebIfc.IFCRELDEFINESBYPROPERTIES);
	for (const relId of propertyRelationships) {
		const rel = webIfcApi.GetLine(modelId, relId) as WebIfc.IFC2X3.IfcRelDefinesByProperties;
		if (!rel.RelatedObjects || !rel.RelatingPropertyDefinition) {
			continue;
		}
		// Process each related object
		for (const relatedObj of rel.RelatedObjects) {
			try {
				const elementId = "value" in relatedObj ? relatedObj.value : relatedObj.expressID;
				const element = webIfcApi.GetLine(modelId, elementId);

				const category = element.constructor.name;
				if (!category) {
					continue;
				}
				if (EXCLUDED_CATEGORIES.has(category)) {
					continue;
				}

				const globalId = element.GlobalId?.value;
				if (!globalId) {
					continue;
				}

				ids.add(globalId);
			} catch (error) {
				log.error(`Error getting element ID: ${error}`);
			}
		}
	}
	return Array.from(ids);
}

/**
 * Create mock LCA data
 * @param ids - The IDs of the elements in the IFC file
 * @returns The mock LCA data
 */
function createMockLcaData(ids: string[]): LcaData[] {
	const lcaData: LcaData[] = [];
	// get a max sequence int < 5

	ids.forEach((id) => {
		const maxSequence = Math.floor(Math.random() * 5);
		for (let i = 0; i < maxSequence; i++) {
			lcaData.push({
				id,
				sequence: i,
				mat_kbob: m(i),
				gwp_relative: n(),
				gwp_absolute: n(),
				penr_relative: n(),
				penr_absolute: n(),
				upb_relative: n(),
				upb_absolute: n(),
			});
		}
	});

	return lcaData;
}

/**
 * Create mock cost data
 * @param ids - The IDs of the elements in the IFC file
 * @returns The mock cost data
 */
function createMockCostData(ids: string[]): CostData[] {
	return ids.map((id) => ({
		id,
		cost: n(),
		cost_unit: n(),
	}));
}

/**
 * Generate a random number
 * @returns A random number between 0 and 100
 */
function n(): number {
	return Math.random() * 100;
}

const MATERIALS = ["concrete", "steel", "wood", "glass", "aluminum", "plastic"];

/**
 * Generate a random material
 * @returns A random material
 */
function m(index: number): string {
	// get a random material from the list
	return MATERIALS[index % MATERIALS.length];
}
