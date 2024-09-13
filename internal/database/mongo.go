package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongo() {
    uri := os.Getenv("MONGO_URI")
    client, err := mongo.NewClient(options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    MongoClient = client
    log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
    dbName := os.Getenv("MONGO_DBNAME")
    return MongoClient.Database(dbName).Collection(collectionName)
}
