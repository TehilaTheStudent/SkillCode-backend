package main

import (
	"fmt"
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
	"go.uber.org/zap"
)

func main() {
	// Ensure the working directory is correct
	utils.EnsureWorkingDirectory()

	// Initialize logger
	logger := config.InitLogger()

	// Load the configuration
	cfg := config.LoadConfigAPI()

	// Initialize MongoDB client
	db.ConnectMongoDB(cfg.MongoDBURI)

	// Initialize the repository and service
	questionRepo := repository.NewQuestionRepository(db.Client.Database(cfg.DBName))
	questionService := service.NewQuestionService(questionRepo)

	// Create a handler with the service instance
	questionHandler := handler.NewQuestionHandler(questionService)

	// Setup Gin router
	r := gin.Default()

	// Add middlewares
	setupMiddlewares(r, logger,cfg.FrontendURL)

	// Register routes
	registerRoutes(r, questionHandler)

	// Start the server
	logger.Info("Starting server", zap.String("port", cfg.Port))
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}
}

// setupMiddlewares adds required middlewares to the Gin router
func setupMiddlewares(r *gin.Engine, logger *zap.Logger, frontendURL string) {
	// Inject logger into Gin context
	r.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	})

	// Add CORS middleware with custom logic to allow the frontend
	corsMiddleware := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			// Match the exact frontend URL
			return origin == frontendURL
		},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		MaxAge:           int(12 * time.Hour / time.Second),                   // Cache preflight request for 12 hours
	})
	r.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	})
}



// registerRoutes registers all application routes
func registerRoutes(r *gin.Engine, questionHandler *handler.QuestionHandler) {
	
	handler.RegisterQuestionRoutes(r, questionHandler)
    handler.RegisterCodeRoutes(r)
	handler.RegisterConfigRoutes(r)
}
