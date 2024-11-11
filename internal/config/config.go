package config

import (
	"log"
	"os"
)

// Config holds all dynamic configuration values
type Config struct {
	MongoURI string
	DBName   string
	Port     string
}

// LoadConfig loads the configuration values from environment variables or sets defaults
func LoadConfig() *Config {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
		log.Println("Using default MongoDB URI")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "skillcode_db"
		log.Println("Using default database name")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("Using default port")
	}

	return &Config{
		MongoURI: mongoURI,
		DBName:   dbName,
		Port:     port,
	}
}
