package main

import (
	"context"
	"fmt"
	"os"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/dependencies"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/handler"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/middleware"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/service"
	tester "github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	fmt.Println("Successfully connected to API server. Namespaces:")

	kubeconfigPath := "/root/.kube/config" // Path to kubeconfig inside container
	apiServerURL := "https://host.docker.internal:37000"



	config8, err := clientcmd.BuildConfigFromFlags(apiServerURL, kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to build config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config8)
	if err != nil {
		fmt.Printf("Failed to create Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Test by listing namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to connect to API server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to API server. Namespaces:")
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
	// Initialize the logger
	logger, err := config.InitLogger()
	if err != nil {
		os.Exit(1)
	}
	// Load the configuration
	err = config.InitGlobalConfigs()
	if err != nil {
		logger.Fatal("Failed to initialize configuration", zap.Error(err))
	}
	// Setup all dependencies
	mongoClient, sharedTester, err := dependencies.SetupAllDependencies()
	if err != nil {
		logger.Fatal("Failed to setup dependencies", zap.Error(err))
	}
	// Initialize handlers
	questionHandler := initializeHandlers(mongoClient, sharedTester)

	// Setup the router with middlewares and routes
	r := setupRouter(logger, questionHandler)

	// Start the server
	logger.Info("Starting server on port", zap.String("port", config.GlobalConfigAPI.Port))
	if err := startServer(r); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}
}

// initializeHandlers sets up the handlers for the application (repository<-service<-handler)
// this is the dependency injection
func initializeHandlers(client *mongo.Client, sharedTester *tester.SharedTester) *handler.QuestionHandler {
	questionRepo := repository.NewQuestionRepository(client.Database(config.GlobalConfigAPI.DBName))
	questionService := service.NewQuestionService(questionRepo, sharedTester)
	return handler.NewQuestionHandler(questionService)
}

// setupRouter configures the router with middlewares and routes fron ; questions, code, config
func setupRouter(logger *zap.Logger, questionHandler *handler.QuestionHandler) *gin.Engine {
	r := gin.Default() //logs every request to the terminal
	middleware.SetupMiddlewares(r, logger, config.GlobalConfigAPI.FrontendURLS)
	handler.RegisterRoutes(r, questionHandler)
	return r
}

// startServer starts the HTTP server
func startServer(r *gin.Engine) error {
	return r.Run(fmt.Sprintf(":%s", config.GlobalConfigAPI.Port))
}
