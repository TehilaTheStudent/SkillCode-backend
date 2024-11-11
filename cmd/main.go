package main

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/db"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handlers"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// MongoDB connection
	cfg := config.LoadConfig() // Load the configuration
	db.ConnectMongoDB(cfg.MongoURI) //initialize MongoDB client

	// Initialize the repository and service
	questionRepo := repository.NewQuestionRepository(db.Client.Database(cfg.DBName))
	questionService := service.NewQuestionService(questionRepo)

	// Create a handler with the service instance
	questionHandler := handlers.NewQuestionHandler(questionService)

	// Gin setup
	r := gin.Default()
	handlers.RegisterQuestionRoutes(r, questionHandler)
	r.Run() // listen and serve on 0.0.0.0:8080
}
