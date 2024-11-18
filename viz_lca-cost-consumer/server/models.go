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
	Id         string  `json:"id"`
	Category   string  `json:"category"`
	CO2e       float32 `json:"co2e"`
	GreyEnergy float32 `json:"greyEnergy"`
	UBP        float32 `json:"UBP"`
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
}
