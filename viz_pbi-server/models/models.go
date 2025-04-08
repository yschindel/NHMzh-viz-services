package models

import "mime/multipart"

// EavElementDataItem represents a single EAV data item
type EavElementDataItem struct {
	Project           string   `json:"project"`             // Project identifier
	Filename          string   `json:"filename"`            // Filename (without path)
	Timestamp         string   `json:"timestamp"`           // Message timestamp
	Id                string   `json:"id"`                  // Unique identifier for the item
	ParamName         string   `json:"param_name"`          // The name of the parameter
	ParamValueString  *string  `json:"param_value_string"`  // The value of the parameter
	ParamValueNumber  *float32 `json:"param_value_number"`  // The value of the parameter
	ParamValueBoolean *bool    `json:"param_value_boolean"` // The value of the parameter
	ParamValueDate    *string  `json:"param_value_date"`    // The value of the parameter
	ParamType         string   `json:"param_type"`          // The type of the parameter (currently always "number")
}

// EavMaterialDataItem represents a single EAV data item
type EavMaterialDataItem struct {
	Project           string   `json:"project"`             // Project identifier
	Filename          string   `json:"filename"`            // Filename (without path)
	Timestamp         string   `json:"timestamp"`           // Message timestamp
	Id                string   `json:"id"`                  // Unique identifier for the item
	Sequence          int      `json:"sequence"`            // The sequence number of the item
	ParamName         string   `json:"param_name"`          // The name of the parameter
	ParamValueString  *string  `json:"param_value_string"`  // The value of the parameter
	ParamValueNumber  *float32 `json:"param_value_number"`  // The value of the parameter
	ParamValueBoolean *bool    `json:"param_value_boolean"` // The value of the parameter
	ParamValueDate    *string  `json:"param_value_date"`    // The value of the parameter
	ParamType         string   `json:"param_type"`          // The type of the parameter (currently always "number")
}

type BlobData struct {
	Container         string         // Container name
	Project           string         // Project identifier
	Filename          string         // Filename (without path)
	Timestamp         string         // Message timestamp
	StorageServiceURL string         // The storage URL of the blob
	BlobID            string         // Internal field containing "project/filename" as a unique identifier
	Blob              multipart.File // The blob data
}
