package data

import (
	cryptrand "crypto/rand"
	"encoding/base64"
	"math/rand"
)

type DataItem struct {
	Id string `bson:"_id"`
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

var categories = [5]string{"walls", "floors", "ceilings", "doors", "windows"}

func randCategory() string {
	return categories[rand.Intn(len(categories))]
}

var epkph = [10]string{"B.07.01", "C.01.01", "C.02.01", "C.03.01", "C.04.01", "C.05.01", "E.02.01", "E.03.01", "G.01.01", "G.02.01"}

func randEbkph() string {
	return epkph[rand.Intn(len(epkph))]
}

var levels = [3]string{"Level 1", "Level 2", "Level 3"}

func randLevel() string {
	return levels[rand.Intn(len(levels))]
}

var fireRatings = [2]string{"A", "B"}

func randFireRating() string {
	return fireRatings[rand.Intn(len(fireRatings))]
}

func newDataItems(numElements int) []DataItem {
	dataItems := make([]DataItem, numElements)
	for i := 0; i < numElements; i++ {
		id, err := generate22CharGUID()
		if err != nil {
			panic(err)
		}
		dataItems[i] = DataItem{
			Id: id,
		}
	}
	return dataItems
}

// this is used to make lca and cost data have the same ids.
// That way the data can be merged into the same item in mongoDB
var DataItems = newDataItems(3)
