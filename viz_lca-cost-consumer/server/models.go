package server

type CustomMinioCredentials struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

type LcaMessage struct {
	Project   string        `json:"project"`
	Filename  string        `json:"filename"`
	Timestamp string        `json:"timestamp"`
	Data      []LcaDataItem `json:"data"`
}

type LcaDataItem struct {
	Id           string  `json:"id"`
	Category     string  `json:"category"`
	GwpAbsolute  float32 `json:"gwp_absolute"`
	GwpRelative  float32 `json:"gwp_relative"`
	PenrAbsolute float32 `json:"penr_absolute"`
	PenrRelative float32 `json:"penr_relative"`
	UbpAbsolute  float32 `json:"ubp_absolute"`
	UbpRelative  float32 `json:"ubp_relative"`
}

type CostMessage struct {
	Project   string         `json:"project"`
	Filename  string         `json:"filename"`
	Timestamp string         `json:"timestamp"`
	Data      []CostDataItem `json:"data"`
}

type CostDataItem struct {
	Id       string  `json:"id"`
	Category string  `json:"category"`
	Cost     float32 `json:"cost"`
	CostUnit float32 `json:"cost_unit"`
}
