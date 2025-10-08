package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connection() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("unable to find .env file", err)
	}

	uri := os.Getenv("MONGODB_URI")

	// Defines the options for the MongoDB client
	opts := options.Client().ApplyURI(uri)

	// Creates a new client and connects to the server
	client, err := mongo.Connect(opts)

	if err != nil {
		log.Fatal("Failed to connect Mongo", err)
	}

	return client
}

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("unable to find .env file", err)
	}

	dbname := os.Getenv("DATABASE_NAME")

	collection := client.Database(dbname).Collection(collectionName)

	fmt.Printf("Opened %s collection\n", collectionName)

	return collection
}
