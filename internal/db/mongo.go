package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() (*mongo.Database, error) {
	uri := "mongodb://mongodb:27017" // must match docker-compose service name
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect mongo: %w", err)
	}

	return client.Database("banking"), nil
}
