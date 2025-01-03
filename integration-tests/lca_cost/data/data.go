package data

import (
	cryptrand "crypto/rand"
	"encoding/base64"
	"math/rand"
)

type DataItem struct {
	Id       string `bson:"_id"`
	Category string `bson:"ebkph"`
}

func generate22CharGUID() (string, error) {
	uuid := make([]byte, 16) // UUID is 16 bytes (128 bits)

	// Generate random bytes
	_, err := cryptrand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version and variant bits for the UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is RFC4122

	// Encode to base64 URL-safe and remove the trailing '=='
	encoded := base64.RawURLEncoding.EncodeToString(uuid)

	return encoded, nil
}

// randFloat generates a random float number between min and max
func randFloat() float32 {
	return rand.Float32()*(100-5) + 5
}

// randCat picks a random category from a predefined list
func randCat() string {
	return categories[rand.Intn(len(categories))]
}

var categories = [10]string{"B07", "C01", "C02", "C03", "C04", "C05", "E02", "E03", "G01", "G02"}

func newDataItems(numElements int) []DataItem {
	dataItems := make([]DataItem, numElements)
	for i := 0; i < numElements; i++ {
		id, err := generate22CharGUID()
		if err != nil {
			panic(err)
		}
		category := randCat()
		dataItems[i] = DataItem{
			Id:       id,
			Category: category,
		}
	}
	return dataItems
}

// this is used to make lca and cost data have the same ids.
// That way the data can be merged into the same item in mongoDB
var DataItems = newDataItems(3)
