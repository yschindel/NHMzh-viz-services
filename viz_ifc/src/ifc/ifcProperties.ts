import * as WebIfc from "web-ifc";
import { log } from "../utils/logger";
import fs from "fs";
import { sendDataToStorage } from "../storage";
import { IFCData, FilePropertyData, PropertyData } from "../types";
import { toEavElementDataItems } from "../data";
import { getEnv } from "../utils/env";

/**
 * Processes an IFC file to extract all properties
 * @param ifcData - The IFC file data
 * @param wasmPath - Path to the web-ifc WASM file
 * @returns Promise with the extracted properties
 */
export async function processIfcProperties(ifcData: IFCData, wasmPath: string): Promise<void> {
	log.info(`Extracting properties from IFC file ${ifcData.filename}`);

	try {
		// Load the IFC file
		const dataArray = new Uint8Array(ifcData.file);
		const webIfcApi = new WebIfc.IfcAPI();
		await webIfcApi.Init();
		webIfcApi.SetWasmPath(wasmPath, true);

		// Open the IFC file
		const modelId = webIfcApi.OpenModel(dataArray);
		log.info("Starting property extraction...");

		// Get all properties
		const startTime = Date.now();
		const elementPropertiessMap = getProperties(webIfcApi, modelId);
		const endTime = Date.now();
		log.debug(`Property extraction completed in ${endTime - startTime}ms`);

		// Close the IFC file
		webIfcApi.CloseModel(modelId);

		// Convert Map to array for final output
		const properties = Array.from(elementPropertiessMap.values());

		// Save properties to a file if we are in dev mode
		if (process.env.NODE_ENV === "development") {
			const jsonData = JSON.stringify(properties, null, 2);
			fs.writeFileSync(`${ifcData.filename}.properties.json`, jsonData);
			log.debug(`Properties saved to ${ifcData.filename}.properties.json`);
		}

		const filePropertyData: FilePropertyData = {
			project: ifcData.project,
			filename: ifcData.filename,
			timestamp: ifcData.timestamp,
			items: properties,
		};

		// use the regular env loading here because it can be empty and that's fine
		const propsToInclude = process.env["IFC_PROPERTIES_TO_INCLUDE"]?.split(",") || [];
		const eavData = toEavElementDataItems(filePropertyData, propsToInclude);

		// Send data to storage
		await sendDataToStorage(eavData);

		log.info(`Successfully extracted properties from ${ifcData.filename}`);
	} catch (error: any) {
		log.error(`Error extracting properties from IFC file: ${error.message}`);
		if (error.stack) {
			log.error(`Error stack: ${error.stack}`);
		}
		throw error;
	}
}

/**
 * Extracts properties from an IFC file
 * Uses a map of expressID to properties for each element for efficient processing
 * @param webIfcApi - The WebIfc API instance
 * @param modelId - The ID of the model to extract properties from
 * @returns A Map of element IDs to their properties
 */
function getProperties(webIfcApi: WebIfc.IfcAPI, modelId: number): Map<number, any> {
	// First, get all property relationships
	const propertyRelationships = webIfcApi.GetLineIDsWithType(modelId, WebIfc.IFCRELDEFINESBYPROPERTIES);
	log.info(`Found ${propertyRelationships.size()} property relationships`);

	// Create a map to store elements and their properties
	const elementsMap = new Map<number, any>();

	// Process each property relationship
	for (let i = 0; i < propertyRelationships.size(); i++) {
		const relId = propertyRelationships.get(i);
		try {
			const rel = webIfcApi.GetLine(modelId, relId) as WebIfc.IFC2X3.IfcRelDefinesByProperties;

			if (!rel.RelatedObjects || !rel.RelatingPropertyDefinition) {
				continue;
			}

			// Get the property set
			const propertySetId = (rel.RelatingPropertyDefinition as WebIfc.Handle<WebIfc.IFC2X3.IfcPropertySetDefinition>).value;
			const propertySet = webIfcApi.GetLine(modelId, propertySetId, true) as WebIfc.IFC2X3.IfcPropertySet;

			if (!propertySet || !propertySet.Name?.value) {
				continue;
			}

			// Process each related object
			for (const relatedObj of rel.RelatedObjects) {
				try {
					const elementId = "value" in relatedObj ? relatedObj.value : relatedObj.expressID;
					const element = webIfcApi.GetLine(modelId, elementId);
					const elementType = element.constructor.name;

					// Find or create element entry
					let elementEntry: PropertyData = elementsMap.get(elementId);
					if (!elementEntry) {
						const level = getElementLevel(webIfcApi, modelId, elementId);
						elementEntry = {
							expressId: elementId,
							globalId: element.GlobalId?.value,
							properties: {
								category: elementType,
								level: level,
							} as Record<string, string | number | boolean>,
						};
						elementsMap.set(elementId, elementEntry);
					}

					// Process properties in the property set
					if (!propertySet.HasProperties) continue;

					for (const prop of propertySet.HasProperties) {
						try {
							const property = "value" in prop ? webIfcApi.GetLine(modelId, prop.value) : prop;

							if (property.type === WebIfc.IFCPROPERTYSINGLEVALUE) {
								const singleValue = property as WebIfc.IFC2X3.IfcPropertySingleValue;
								const propertyName = singleValue.Name?.value;
								let propertyValue = singleValue.NominalValue?.value;
								// check if the property value is an array
								if (Array.isArray(propertyValue)) {
									propertyValue = propertyValue.join(", ");
								}
								if (propertyName && propertyValue !== undefined) {
									elementEntry.properties[propertyName] = propertyValue;
								}
							} else if (property.type === WebIfc.IFCPROPERTYENUMERATEDVALUE) {
								const enumValue = property as WebIfc.IFC2X3.IfcPropertyEnumeratedValue;
								const propertyName = enumValue.Name?.value;
								const propertyValue = enumValue.EnumerationValues?.map((v) => v.value).join(", ");
								if (propertyName && propertyValue) {
									elementEntry.properties[propertyName] = propertyValue;
								}
							} else if (property.type === WebIfc.IFCPROPERTYBOUNDEDVALUE) {
								const boundedValue = property as WebIfc.IFC2X3.IfcPropertyBoundedValue;
								const propertyName = boundedValue.Name?.value;
								const lowerBound = boundedValue.LowerBoundValue?.value;
								const upperBound = boundedValue.UpperBoundValue?.value;
								if (propertyName && (lowerBound !== undefined || upperBound !== undefined)) {
									elementEntry.properties[propertyName] = `${lowerBound ?? ""} - ${upperBound ?? ""}`;
								}
							} else if (property.type === WebIfc.IFCPROPERTYLISTVALUE) {
								const listValue = property as WebIfc.IFC2X3.IfcPropertyListValue;
								const propertyName = listValue.Name?.value;
								const propertyValue = listValue.ListValues?.map((v) => v.value).join(", ");
								if (propertyName && propertyValue) {
									elementEntry.properties[propertyName] = propertyValue;
								}
							} else {
								// Skip unknown property types
								continue;
							}
						} catch (error) {
							log.warn(`Error processing property for element ${elementId}: ${error}`);
						}
					}
				} catch (error) {
					log.warn(`Error processing related object in relationship ${rel.expressID}: ${error}`);
				}
			}
		} catch (error) {
			log.warn(`Error processing property relationship ${relId}: ${error}`);
		}
	}

	return elementsMap;
}

function getElementLevel(webIfcApi: WebIfc.IfcAPI, modelId: number, elementId: number): string | undefined {
	try {
		// Get all spatial containment relationships
		const spatialRelations = webIfcApi.GetLineIDsWithType(modelId, WebIfc.IFCRELCONTAINEDINSPATIALSTRUCTURE);

		for (let i = 0; i < spatialRelations.size(); i++) {
			const relId = spatialRelations.get(i);
			const rel = webIfcApi.GetLine(modelId, relId) as WebIfc.IFC2X3.IfcRelContainedInSpatialStructure;

			// Check if our element is in this relationship
			if (rel.RelatedElements.some((el) => ("value" in el ? el.value : el.expressID) === elementId)) {
				// Get the spatial structure (should be a level)
				const spatialStructureId = "value" in rel.RelatingStructure ? rel.RelatingStructure.value : rel.RelatingStructure.expressID;
				const spatialStructure = webIfcApi.GetLine(modelId, spatialStructureId);
				if (spatialStructure.type === WebIfc.IFCBUILDINGSTOREY) {
					return spatialStructure.Name?.value || spatialStructure.GlobalId?.value || "No Level";
				}
			}
		}
	} catch (error) {
		log.warn(`Error getting level for element ${elementId}: ${error}`);
	}
	return "No Level";
}
