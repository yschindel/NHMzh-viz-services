package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Element struct {
	Id         string  `bson:"_id"`
	Timestamp  string  `bson:"timestamp"`
	Category   string  `bson:"category"`
	Cost       float32 `bson:"cost"`
	CO2e       float32 `json:"co2e"`
	GreyEnergy float32 `json:"greyEnergy"`
	UBP        float32 `json:"UBP"`
}

type Reader struct {
	client *mongo.Client
}

func NewReader(uri string) *Reader {
	log.Printf("connecting to MongoDB: %s", uri)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	} else {
		log.Printf("Connected to MongoDB")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	} else {
		log.Printf("Ping to MongoDB successful")
	}

	return &Reader{
		client: client,
	}
}

func (r *Reader) ReadAllElements(dbName, collectionName string) ([]Element, error) {
	collection := r.client.Database(dbName).Collection(collectionName)

	// Find all documents in the collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(context.TODO())

	var elements []Element
	for cursor.Next(context.TODO()) {
		var element Element
		if err := cursor.Decode(&element); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		elements = append(elements, element)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return elements, nil
}
