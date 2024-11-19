package dependencies

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var Client *mongo.Client
var IsMongoDBHealthy = true
// ConnectMongoDB initializes the MongoDB client with the provided URI and performs a health check.
func ConnectMongoDB(uri string, logger *zap.Logger) {
	// Set a timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client options and connect
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err)) // Fail fast if connection fails
	}

	// Perform a health check (Ping) to verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.Error(err)) // Fail fast if Ping fails
	}

	logger.Info("Connected to MongoDB successfully!")
	Client = client
}

// HealthCheck verifies the MongoDB connection by performing a Ping.
func HealthCheck(logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := Client.Ping(ctx, nil); err != nil {
		logger.Error("MongoDB health check failed",
			zap.Error(err),
			zap.String("service", "MongoDB"),
		)
		IsMongoDBHealthy=false
		return err
	}
	IsMongoDBHealthy = true // Update global health status
	return nil
}

func StartBackgroundHealthCheck(logger *zap.Logger, interval time.Duration) {
	go func() {
		for {
			err := HealthCheck(logger)
			if err != nil {
				logger.Error("MongoDB health check failed", zap.Error(err))
			}
			time.Sleep(interval)
		}
	}()
}
