package server

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Writer struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewWriter(uri, dbName, collectionName string) *Writer {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &Writer{
		client:     client,
		collection: collection,
	}
}

func (w *Writer) InsertDocument(document interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := w.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result, nil
}
