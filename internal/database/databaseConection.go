package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongodb := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

var client *mongo.Client = DBinstance()

func OpenConnection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("").Collection(collectionName)
}
