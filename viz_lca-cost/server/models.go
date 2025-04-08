package server

// LcaMessage represents a message with Life Cycle Assessment data
type LcaMessage struct {
	Project   string    `json:"project"`   // Project identifier
	Filename  string    `json:"filename"`  // Filename (without path)
	Timestamp string    `json:"timestamp"` // Message timestamp
	Data      []LcaItem `json:"data"`      // Array of LCA data items
	FileID    string    `json:"-"`         // Internal field containing "project/filename" as a unique identifier
}

// LcaItem represents a single LCA data item
type LcaItem struct {
	Id           string  `json:"id"`            // Unique identifier for the item
	Sequence     int     `json:"sequence"`      // The sequence number of the material in the element (STARTING FROM 0)
	MaterialKbob string  `json:"mat_kbob"`      // KBOB material reference
	GwpAbsolute  float32 `json:"gwp_absolute"`  // Global Warming Potential (absolute)
	GwpRelative  float32 `json:"gwp_relative"`  // Global Warming Potential (relative)
	PenrAbsolute float32 `json:"penr_absolute"` // Primary Energy Non-Renewable (absolute)
	PenrRelative float32 `json:"penr_relative"` // Primary Energy Non-Renewable (relative)
	UbpAbsolute  float32 `json:"ubp_absolute"`  // Environmental Impact Points (absolute)
	UbpRelative  float32 `json:"ubp_relative"`  // Environmental Impact Points (relative)
}

// CostMessage represents a message with cost data
type CostMessage struct {
	Project   string     `json:"project"`   // Project identifier
	Filename  string     `json:"filename"`  // Filename (without path)
	Timestamp string     `json:"timestamp"` // Message timestamp
	Data      []CostItem `json:"data"`      // Array of cost data items
	FileID    string     `json:"-"`         // Internal field containing "project/filename" as a unique identifier
}

// CostItem represents a single cost data item
type CostItem struct {
	Id       string  `json:"id"`        // Unique identifier for the item
	Cost     float32 `json:"cost"`      // Cost value
	CostUnit float32 `json:"cost_unit"` // Cost unit
}

// EavElementRow represents a single EAV data item (in this case it's used for cost data)
type EavElementRow struct {
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

// EavMaterialRow represents a single EAV data item
type EavMaterialRow struct {
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
