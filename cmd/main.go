package main

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/db"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handlers"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	// MongoDB connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default MongoDB URI
	}
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "skillcode_db" // fallback if MONGO_DB is not set
	}
	db.ConnectMongoDB(mongoURI) //initialize MongoDB client

	// Initialize the repository and service
	questionRepo := repository.NewQuestionRepository(db.Client.Database(dbName))
	questionService := service.NewQuestionService(questionRepo)

	// Create a handler with the service instance
	questionHandler := handlers.NewQuestionHandler(questionService)

	// Gin setup
	r := gin.Default()
	handlers.RegisterQuestionRoutes(r, questionHandler)
	r.Run() // listen and serve on 0.0.0.0:8080
}
