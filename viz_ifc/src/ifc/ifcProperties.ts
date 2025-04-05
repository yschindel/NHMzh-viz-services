import * as WebIfc from "web-ifc";
import { log } from "../utils/logger";
import fs from "fs";

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

		const startTime = Date.now();

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

				const psetName = propertySet.Name.value;

				// Process each related object
				for (const relatedObj of rel.RelatedObjects) {
					try {
						const elementId = "value" in relatedObj ? relatedObj.value : relatedObj.expressID;
						const element = webIfcApi.GetLine(modelId, elementId);
						const elementType = element.constructor.name;

						// Find or create element entry
						let elementEntry = elementsMap.get(elementId);
						if (!elementEntry) {
							elementEntry = {
								expressId: elementId,
								type: elementType,
								properties: {},
							};
							if (element.GlobalId?.value) {
								elementEntry.globalId = element.GlobalId.value;
							}
							elementsMap.set(elementId, elementEntry);
						}

						// Initialize property set if it doesn't exist
						if (!elementEntry.properties[psetName]) {
							elementEntry.properties[psetName] = {};
						}

						// Process properties in the property set
						if (propertySet.HasProperties) {
							for (const prop of propertySet.HasProperties) {
								try {
									const property = "value" in prop ? webIfcApi.GetLine(modelId, prop.value) : prop;

									if (property.type === WebIfc.IFCPROPERTYSINGLEVALUE) {
										const singleValue = property as WebIfc.IFC2X3.IfcPropertySingleValue;
										const propertyName = singleValue.Name?.value;
										const propertyValue = singleValue.NominalValue?.value;
										if (propertyName && propertyValue !== undefined) {
											elementEntry.properties[psetName][propertyName] = propertyValue;
										}
									} else if (property.type === WebIfc.IFCPROPERTYENUMERATEDVALUE) {
										const enumValue = property as WebIfc.IFC2X3.IfcPropertyEnumeratedValue;
										const propertyName = enumValue.Name?.value;
										const propertyValue = enumValue.EnumerationValues?.map((v) => v.value).join(", ");
										if (propertyName && propertyValue) {
											elementEntry.properties[psetName][propertyName] = propertyValue;
										}
									} else if (property.type === WebIfc.IFCPROPERTYBOUNDEDVALUE) {
										const boundedValue = property as WebIfc.IFC2X3.IfcPropertyBoundedValue;
										const propertyName = boundedValue.Name?.value;
										const lowerBound = boundedValue.LowerBoundValue?.value;
										const upperBound = boundedValue.UpperBoundValue?.value;
										if (propertyName && (lowerBound !== undefined || upperBound !== undefined)) {
											elementEntry.properties[psetName][propertyName] = {
												lowerBound,
												upperBound,
											};
										}
									} else if (property.type === WebIfc.IFCPROPERTYLISTVALUE) {
										const listValue = property as WebIfc.IFC2X3.IfcPropertyListValue;
										const propertyName = listValue.Name?.value;
										const propertyValues = listValue.ListValues?.map((v) => v.value);
										if (propertyName && propertyValues?.length) {
											elementEntry.properties[psetName][propertyName] = propertyValues;
										}
									}
								} catch (error) {
									log.warn(`Error processing property for element ${elementId}: ${error}`);
								}
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

		const endTime = Date.now();
		log.debug(`Property extraction completed in ${endTime - startTime}ms`);

		webIfcApi.CloseModel(modelId);

		// Convert Map to array for final output
		const elementsWithProperties = Array.from(elementsMap.values());

		// Save properties to a file
		const jsonData = JSON.stringify(elementsWithProperties, null, 2);
		fs.writeFileSync(`${ifcData.Filename}.properties.json`, jsonData);

		log.info(`Successfully extracted properties from ${ifcData.Filename}`);
		log.debug(`Found ${elementsWithProperties.length} elements with properties`);
	} catch (error: any) {
		log.error(`Error extracting properties from IFC file: ${error.message}`);
		if (error.stack) {
			log.error(`Error stack: ${error.stack}`);
		}
		throw error;
	}
}
