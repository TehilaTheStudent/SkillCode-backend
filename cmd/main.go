package main

import (
	"log"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/db"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handler"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Ensure the working directory is correct
	utils.EnsureWorkingDirectory()

	// Load the configuration
	cfg := config.LoadConfig()

	// Initialize MongoDB client
	db.ConnectMongoDB(cfg.MongoDBURI)

	// Initialize the repository and service
	questionRepo := repository.NewQuestionRepository(db.Client.Database(cfg.DBName))
	questionService := service.NewQuestionService(questionRepo)

	// Create a handler with the service instance
	questionHandler := handler.NewQuestionHandler(questionService)

	// Setup Gin router
	r := gin.Default()

	// Register routes
	handler.RegisterQuestionRoutes(r, questionHandler)

	// Run the server
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
