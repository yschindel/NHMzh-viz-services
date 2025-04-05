// Interface to carry data about the file to be uploaded to storage
// Casing does not matter if name only differs by casing when parsing in GO
interface FileData {
	file: Buffer;
	project: string;
	filename: string;
	timestamp: string;
	fileId: string;
}

// Interface to carry data about the IFC file
// Casing does not matter if name only differs by casing when parsing in GO
interface IFCData {
	file: Buffer;
	project: string;
	filename: string;
	timestamp: string;
	fileId: string;
}

interface PropertyData {
	expressId: number;
	globalId: string;
	properties: Record<string, string | number | boolean>;
}

interface FilePropertyData {
	project: string;
	filename: string;
	timestamp: string;
	items: PropertyData[];
}

// camel case to match sql column names
interface ElementDataEAV {
	project: string;
	filename: string;
	timestamp: string;
	id: string;
	param_name: string;
	param_value_string: string | undefined;
	param_value_number: number | undefined;
	param_value_boolean: boolean | undefined;
	param_value_date: string | undefined;
	param_type: string;
}

export type { FileData, IFCData, FilePropertyData, ElementDataEAV, PropertyData };
