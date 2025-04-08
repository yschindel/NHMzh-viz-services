import { FilePropertyData, ElementDataEAV } from "./types";
import { log } from "./utils/logger";
/**
 * Converts a FilePropertyData object to an array of ElementDataEAV objects
 * @param data - The FilePropertyData object to convert
 * @param include - Optional array of parameter names to include. If empty, all parameters will be included.
 * @returns An array of ElementDataEAV objects
 */
export function toEavElementDataItems(data: FilePropertyData, include: string[] = []): ElementDataEAV[] {
	const eavItems: ElementDataEAV[] = [];
	for (const item of data.items) {
		if (!item.properties || Object.keys(item.properties).length === 0 || !item.globalId) continue;
		for (const [paramName, paramValue] of Object.entries(item.properties)) {
			if (!paramName || !paramValue) continue;
			if (include.length > 0 && !include.includes(paramName)) continue;

			const eavItem: ElementDataEAV = {
				project: data.project,
				filename: data.filename,
				timestamp: data.timestamp,
				id: item.globalId,
				param_name: paramName,
				param_value_string: undefined,
				param_value_number: undefined,
				param_value_boolean: undefined,
				param_value_date: undefined,
				param_type: "string",
			};

			if (typeof paramValue === "string") {
				eavItem.param_value_string = paramValue;
				eavItem.param_type = "string";
			} else if (typeof paramValue === "number") {
				eavItem.param_value_number = paramValue;
				eavItem.param_type = "number";
			} else if (typeof paramValue === "boolean") {
				eavItem.param_value_boolean = paramValue;
				eavItem.param_type = "boolean";
			} else {
				eavItem.param_value_string = JSON.stringify(paramValue);
				eavItem.param_type = "string";
			}

			eavItems.push(eavItem);
		}
	}

	return eavItems;
}
