package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() *mongo.Client {

	mongoDB := os.Getenv("MONGO_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDB))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully connected to  MongoDB server")
	return client
}

var Client *mongo.Client = InitDB()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Restaurant_Management").Collection(collectionName)
	return collection
}
