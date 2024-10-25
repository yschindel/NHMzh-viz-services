package server

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Element struct {
	DataItem
	Timestamp string `bson:"timestamp"`
}

type Writer struct {
	client *mongo.Client
}

func NewWriter(uri string) *Writer {
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

	return &Writer{
		client: client,
	}
}

func (w *Writer) UpsertElements(db, col string, els []Element) error {
	collection := w.client.Database(db).Collection(col)

	// Prepare the bulk operations
	var operations []mongo.WriteModel
	for _, el := range els {
		filter := bson.M{"_id": el.Id}
		update := bson.M{
			"$set": bson.M{
				"timestamp":  el.Timestamp,
				"category":   el.Category,
				"co2e":       el.CO2e,
				"greyEnergy": el.GreyEnergy,
				"UBP":        el.UBP,
			},
		}
		upsert := true
		// Create an upsert operation for each element
		operation := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(upsert)
		operations = append(operations, operation)
	}

	// Execute the bulk write operation
	bulkOpts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(context.TODO(), operations, bulkOpts)
	if err != nil {
		log.Fatal("Failed to execute bulkWrite: ", err)
		return err
	}

	fmt.Printf("BulkWrite completed with %d upserts and %d updates\n", result.UpsertedCount, result.ModifiedCount)

	return nil
}

func UnpackMessage(message LcaMessage) []Element {
	elements := make([]Element, len(message.Data))
	for i, dataItem := range message.Data {
		elements[i] = Element{
			DataItem:  dataItem,
			Timestamp: message.Timestamp,
		}
	}
	return elements
}
