package repositories

import (
	"context"
	"log"
	"time"

	"cms-server/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertContent(collection *mongo.Collection, content models.Content) (*mongo.InsertOneResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, content)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    return result, nil
}

// Other CRUD operations like FindContent, UpdateContent, DeleteContent...
