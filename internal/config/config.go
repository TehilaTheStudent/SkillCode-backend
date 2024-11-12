package config

import (
	"log"
	"os"
)

// Config holds all dynamic configuration values
type Config struct {
	MongoDBURI string
	DBName     string
	Port       string
}

// LoadConfig loads the application configuration from environment variables or a config file.
func LoadConfig() *Config {
	return &Config{
		MongoDBURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:     getEnv("MONGO_DB", "skillcode_db"),
		Port:       getEnv("PORT", "8080"),
	}
}

// getEnv retrieves the value of the environment variable named by the key or returns the default value if the variable is not set
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Println("Using default value for", key)
		return defaultValue
	}
	return value
}
