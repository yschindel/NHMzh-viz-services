package mongo

type Element struct {
	Id         string  `bson:"_id"`
	Timestamp  string  `bson:"timestamp"`
	Category   string  `bson:"category"`
	Cost       float32 `bson:"cost"`
	CO2e       float32 `json:"co2e"`
	GreyEnergy float32 `json:"greyEnergy"`
	UBP        float32 `json:"UBP"`
}
