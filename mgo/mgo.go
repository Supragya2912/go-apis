package mgo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	Users  *mongo.Collection
)

func InitMongoDB() {

	uri := "mongodb://localhost:27017"

	// MongoDB client options
	clientOptions := options.Client().ApplyURI(uri)

	log.Println("Connecting to MongoDB")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	Client = client
	Users = client.Database("go-apis").Collection("users")
	log.Println("MongoDB connection established")
}
