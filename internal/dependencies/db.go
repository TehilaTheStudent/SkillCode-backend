package dependencies

import (
	"context"
	"fmt"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var IsMongoDBHealthy = true

// initializeDatabase sets up the database connection and starts health checks
func InitializeDatabase() error {
	if err := ConnectMongoDB(config.GlobalConfigAPI.MongoDBURI); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	StartBackgroundHealthCheck(30 * time.Second)
	return nil
}

// ConnectMongoDB initializes the MongoDB client with the provided URI and performs a health check.
func ConnectMongoDB(uri string) error {
	// Set a timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client options and connect
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Perform a health check (Ping) to verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	Client = client
	return nil
}

// HealthCheck verifies the MongoDB connection by performing a Ping.
func HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := Client.Ping(ctx, nil); err != nil {
		IsMongoDBHealthy = false
		return fmt.Errorf("MongoDB health check failed: %w", err)
	}
	IsMongoDBHealthy = true // Update global health status
	return nil
}

func StartBackgroundHealthCheck(interval time.Duration) {
	go func() {
		for {
			_ = HealthCheck()
			time.Sleep(interval)
		}
	}()
}
