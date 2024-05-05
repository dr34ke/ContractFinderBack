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
var databaseName string = "ContractFinder"

func OpenConnection(collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}

func GetOne[T any](client *mongo.Client, collectionName string, filter interface{}) (*T, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result T
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetMany[T any](collectionName string, filter interface{}) ([]T, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func Insert[T any](collectionName string, document T) (*mongo.InsertOneResult, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Update(collectionName string, filter, update interface{}) (*mongo.UpdateResult, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Aggregate[T any](collectionName string, pipeline mongo.Pipeline) ([]T, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func Count(collectionName string, filter interface{}) (int64, error) {
	collection := client.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}