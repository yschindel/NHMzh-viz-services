// Interface to carry data about the file to be uploaded to storage
interface FileData {
	File: Buffer;
	Project: string;
	Filename: string;
	Timestamp: string;
	FileID: string;
}

// Interface to carry data about the IFC file
interface IFCData {
	File: Buffer;
	Project: string;
	Filename: string;
	Timestamp: string;
	FileID: string;
}
