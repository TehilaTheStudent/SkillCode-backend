package dependencies

import (
	"context"
	"fmt"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



// initializeDatabase sets up the database connection 
func InitializeDatabase() (*mongo.Client, error) {
	client, err := ConnectMongoDB(config.GlobalConfigAPI.MongoDBURI)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return client, nil
}

// ConnectMongoDB initializes the MongoDB client with the provided URI 
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	// Set a timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client options and connect
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil,fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Perform a health check (Ping) to verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil,fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client,nil
}
