import * as WebIfc from "web-ifc";
import { log } from "../utils/logger";
import fs from "fs";
// Helper function to get properties for a single element
async function getElementProperties(webIfcApi: WebIfc.IfcAPI, modelId: number, elementId: number) {
	const element = webIfcApi.GetLine(modelId, elementId);
	const properties: any = {
		expressId: elementId,
		type: element.constructor.name,
		properties: {},
	};

	// Only add GlobalId if it exists
	if (element.GlobalId?.value) {
		properties.globalId = element.GlobalId.value;
	}

	try {
		const propertyRelationships = webIfcApi.GetLineIDsWithType(modelId, WebIfc.IFCRELDEFINESBYPROPERTIES);

		for (let i = 0; i < propertyRelationships.size(); i++) {
			const relId = propertyRelationships.get(i);
			const rel = webIfcApi.GetLine(modelId, relId) as WebIfc.IFC2X3.IfcRelDefinesByProperties;

			// Check if the relationship is related to our element
			const isRelated = rel.RelatedObjects?.some((relObj: any) => relObj.value === elementId) ?? false;

			if (isRelated && rel.RelatingPropertyDefinition) {
				const propertySetId = (rel.RelatingPropertyDefinition as WebIfc.Handle<WebIfc.IFC2X3.IfcPropertySetDefinition>).value;
				const propertySet = webIfcApi.GetLine(modelId, propertySetId, true) as WebIfc.IFC2X3.IfcPropertySet;

				const psetName = propertySet.Name?.value;
				if (psetName) {
					properties.properties[psetName] = {};

					propertySet.HasProperties?.forEach((prop: WebIfc.IFC2X3.IfcProperty | WebIfc.Handle<WebIfc.IFC2X3.IfcProperty>) => {
						try {
							const property = "value" in prop ? webIfcApi.GetLine(modelId, prop.value) : prop;

							// Handle different property types
							if (property.type === WebIfc.IFCPROPERTYSINGLEVALUE) {
								const singleValue = property as WebIfc.IFC2X3.IfcPropertySingleValue;
								const propertyName = singleValue.Name?.value;
								const propertyValue = singleValue.NominalValue?.value;
								if (propertyName && propertyValue !== undefined) {
									properties.properties[psetName][propertyName] = propertyValue;
								}
							} else if (property.type === WebIfc.IFCPROPERTYENUMERATEDVALUE) {
								const enumValue = property as WebIfc.IFC2X3.IfcPropertyEnumeratedValue;
								const propertyName = enumValue.Name?.value;
								const propertyValue = enumValue.EnumerationValues?.map((v) => v.value).join(", ");
								if (propertyName && propertyValue) {
									properties.properties[psetName][propertyName] = propertyValue;
								}
							} else if (property.type === WebIfc.IFCPROPERTYBOUNDEDVALUE) {
								const boundedValue = property as WebIfc.IFC2X3.IfcPropertyBoundedValue;
								const propertyName = boundedValue.Name?.value;
								const lowerBound = boundedValue.LowerBoundValue?.value;
								const upperBound = boundedValue.UpperBoundValue?.value;
								if (propertyName && (lowerBound !== undefined || upperBound !== undefined)) {
									properties.properties[psetName][propertyName] = {
										lowerBound,
										upperBound,
									};
								}
							} else if (property.type === WebIfc.IFCPROPERTYLISTVALUE) {
								const listValue = property as WebIfc.IFC2X3.IfcPropertyListValue;
								const propertyName = listValue.Name?.value;
								const propertyValues = listValue.ListValues?.map((v) => v.value);
								if (propertyName && propertyValues?.length) {
									properties.properties[psetName][propertyName] = propertyValues;
								}
							}
						} catch (error) {
							log.warn(`Error processing property for element ${elementId}: ${error}`);
						}
					});
				}
			}
		}
	} catch (error) {
		log.warn(`Error processing property relationships for element ${elementId}: ${error}`);
	}

	return properties;
}

interface ElementProperties {
	expressId: number;
	globalId: string;
	type: string;
	properties: Record<string, Record<string, any>>;
}

interface ElementsByType {
	[key: string]: ElementProperties[];
}

// Main function to get all elements and their properties
export async function getAllElementsWithProperties(webIfcApi: WebIfc.IfcAPI, modelId: number) {
	const allElementsWithProperties: ElementsByType = {};

	// Get all line IDs in the model
	const allLineIds = webIfcApi.GetAllLines(modelId);

	// Create a map to store unique types and their names
	const typeMap = new Map<number, string>();

	// First pass: collect all unique types
	for (let i = 0; i < allLineIds.size(); i++) {
		const lineId = allLineIds.get(i);
		const line = webIfcApi.GetLine(modelId, lineId);
		if (line && line.type) {
			typeMap.set(line.type, line.constructor.name);
		}
	}

	// Second pass: get elements for each type
	for (const [type, typeName] of typeMap) {
		allElementsWithProperties[typeName] = [];
		const elements = new Set(webIfcApi.GetLineIDsWithType(modelId, type));

		// Get properties for each element of this type
		for (const elementId of elements) {
			try {
				const elementProperties = await getElementProperties(webIfcApi, modelId, elementId);
				allElementsWithProperties[typeName].push(elementProperties);
			} catch (error) {
				console.warn(`Error getting properties for ${typeName} element ${elementId}:`, error);
			}
		}

		console.log(`Processed ${elements.size} ${typeName} elements`);
	}

	return allElementsWithProperties;
}

/**
 * Processes an IFC file to extract all properties
 * @param ifcData - The IFC file data
 * @param wasmPath - Path to the web-ifc WASM file
 * @returns Promise with the extracted properties
 */
export async function processIfcProperties(ifcData: IFCData, wasmPath: string): Promise<void> {
	log.info(`Extracting properties from IFC file ${ifcData.Filename}`);

	try {
		const dataArray = new Uint8Array(ifcData.File);
		const webIfcApi = new WebIfc.IfcAPI();
		await webIfcApi.Init();
		webIfcApi.SetWasmPath(wasmPath, true);

		const modelId = webIfcApi.OpenModel(dataArray);
		log.info("Starting property extraction...");
		const properties = await getAllElementsWithProperties(webIfcApi, modelId);
		webIfcApi.CloseModel(modelId);

		// save properties to a file
		const jsonData = JSON.stringify(properties, null, 2);
		fs.writeFileSync(`${ifcData.Filename}.properties.json`, jsonData);

		log.info(`Successfully extracted properties from ${ifcData.Filename}`);
	} catch (error: any) {
		log.error(`Error extracting properties from IFC file: ${error.message}`);
		if (error.stack) {
			log.error(`Error stack: ${error.stack}`);
		}
		throw error;
	}
}
