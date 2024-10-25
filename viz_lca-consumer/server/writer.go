package server

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Writer struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewWriter(uri, dbName, collectionName string) *Writer {
	log.Printf("connecting to MongoDB: %s, db: %s, collection: %s", uri, dbName, collectionName)

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

	collection := client.Database(dbName).Collection(collectionName)

	return &Writer{
		client:     client,
		collection: collection,
	}
}

func (w *Writer) InsertDoc(message LcaMessage) (*mongo.InsertOneResult, error) {
	document, err := bson.Marshal(message)
	if err != nil {
		log.Printf("could not marshal document: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("inserting document...")
	result, err := w.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	} else {
		log.Printf("document inserted: %s", result.InsertedID)
	}

	return result, nil
}
