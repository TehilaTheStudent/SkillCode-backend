package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/db"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handler"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
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

	// Add CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow Nuxt frontend
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		MaxAge:           int(12 * time.Hour / time.Second),                   // Cache preflight request for 12 hours

	})

	// Wrap Gin router with CORS middleware
	r.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	})

	// Register routes
	handler.RegisterQuestionRoutes(r, questionHandler)

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
