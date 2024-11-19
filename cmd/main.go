package main

import (
	"fmt"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/dependencies"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handler"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/middleware"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Ensure the working directory is correct
	dependencies.EnsureWorkingDirectory()
	// Initialize the logger
	logger := config.InitLogger()
	// Load the configuration
	cfg := config.LoadConfigAPI()

	// Initialize the database connection and start health checks
	initializeDatabase(cfg, logger)
	// Initialize handlers
	questionHandler := initializeHandlers(cfg)

	// Setup the router with middlewares and routes
	r := setupRouter(logger, cfg, questionHandler)

	// Start the server
	startServer(r, logger, cfg)
}

// initializeDatabase sets up the database connection and starts health checks
func initializeDatabase(cfg *config.ConfigAPI, logger *zap.Logger) {
	dependencies.ConnectMongoDB(cfg.MongoDBURI, logger)
	dependencies.StartBackgroundHealthCheck(logger, 30*time.Second)
}

// initializeHandlers sets up the handlers for the application
func initializeHandlers(cfg *config.ConfigAPI) *handler.QuestionHandler {
	questionRepo := repository.NewQuestionRepository(dependencies.Client.Database(cfg.DBName))
	questionService := service.NewQuestionService(questionRepo)
	return handler.NewQuestionHandler(questionService)
}

// setupRouter configures the router with middlewares and routes
func setupRouter(logger *zap.Logger, cfg *config.ConfigAPI, questionHandler *handler.QuestionHandler) *gin.Engine {
	r := gin.Default() //logs every request to the terminal
	middleware.SetupMiddlewares(r, logger, cfg.FrontendURL)
	handler.RegisterRoutes(r, questionHandler)
	return r
}

// startServer starts the HTTP server
func startServer(r *gin.Engine, logger *zap.Logger, cfg *config.ConfigAPI) {
	logger.Info("Starting server", zap.String("port", cfg.Port))
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}
}
