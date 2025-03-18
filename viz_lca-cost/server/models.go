package server

type CustomMinioCredentials struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

// LcaMessage represents a message with Life Cycle Assessment data
type LcaMessage struct {
	Project   string        `json:"project"`   // Project identifier
	Filename  string        `json:"filename"`  // Filename (without path)
	Timestamp string        `json:"timestamp"` // Message timestamp
	Data      []LcaDataItem `json:"data"`      // Array of LCA data items
	FileID    string        `json:"-"`         // Internal field containing "project/filename" as a unique identifier
}

// LcaDataItem represents a single LCA data item
type LcaDataItem struct {
	Id           string  `json:"id"`            // Unique identifier for the item
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
	Project   string         `json:"project"`   // Project identifier
	Filename  string         `json:"filename"`  // Filename (without path)
	Timestamp string         `json:"timestamp"` // Message timestamp
	Data      []CostDataItem `json:"data"`      // Array of cost data items
	FileID    string         `json:"-"`         // Internal field containing "project/filename" as a unique identifier
}

// CostDataItem represents a single cost data item
type CostDataItem struct {
	Id           string  `json:"id"`            // Unique identifier for the item
	Category     string  `json:"category"`      // The element category (e.g. "walls" or "floors")
	Level        string  `json:"level"`         // The level of the element (e.g. "wall" or "floor")
	IsStructural bool    `json:"is_structural"` // Whether the element is structural (e.g. true or false)
	FireRating   string  `json:"fire_rating"`   // The fire rating of the element (e.g. "A" or "B")
	Ebkph        string  `json:"ebkph"`         // eBKP-H category (e.g., "C.02.95")
	Ebkph1       string  `json:"-"`             // First component of eBKP-H (e.g., "C" from "C.02.95")
	Ebkph2       string  `json:"-"`             // Second component of eBKP-H (e.g., "02" from "C.02.95")
	Ebkph3       string  `json:"-"`             // Third component of eBKP-H (e.g., "95" from "C.02.95")
	Cost         float32 `json:"cost"`          // Cost value
	CostUnit     float32 `json:"cost_unit"`     // Cost unit
}
