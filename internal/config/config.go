package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

// Config holds the dynamic configuration values
type ConfigSandbox struct {
	ImageName      string
	DockerFilePath string
	PodFilePath    string
	PodName        string
	ClusterName    string
	UserCodePath   string
}

// NewSandboxConfig creates a new sandbox configuration for a given language.
func NewSandboxConfig(language model.PredefinedSupportedLanguage) (*ConfigSandbox, error) {
	baseDir := "./assets"
	// Dynamically derive language paths
	// Convert language to lowercase
	languageStr := strings.ToLower(string(language))
	langDir := fmt.Sprintf("%s/%s", baseDir, languageStr)
	defaultPodFilePath := fmt.Sprintf("%s/pod.yaml", langDir)
	clusterName := "my-cluster"

	dockerFilePath := fmt.Sprintf("%s/Dockerfile", langDir)
	userCodePath := fmt.Sprintf("%s/run_tests.%s", langDir, getFileExtension(language))

	// Validate supported language
	if getFileExtension(language) == "" {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	// Create sandbox configuration
	config := &ConfigSandbox{
		ImageName:      fmt.Sprintf("%s-test:latest", languageStr),
		DockerFilePath: dockerFilePath,
		PodFilePath:    defaultPodFilePath,
		PodName:        "pod",
		ClusterName:    clusterName,
		UserCodePath:   userCodePath,
	}

	return config, nil
}

// getFileExtension returns the file extension for the given language.
func getFileExtension(language model.PredefinedSupportedLanguage) string {
	switch language {
	case model.Python:
		return "py"
	case model.JavaScript:
		return "js"
	case model.Java:
		return "java"
	case model.Go:
		return "go"
	case model.CSharp:
		return "cs"
	case model.Cpp:
		return "cpp"
	default:
		return "" // Unsupported language
	}
}

// Config holds all dynamic configuration values
type ConfigAPI struct {
	MongoDBURI  string
	DBName      string
	Port        string
	FrontendURL string
	Base        string
}

// LoadConfigAPI loads the application configuration from environment variables or a config file.
func LoadConfigAPI() *ConfigAPI {
	return &ConfigAPI{
		MongoDBURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:      getEnv("MONGO_DB", "skillcode_db"),
		Port:        getEnv("PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "*"),
		Base:        "skillcode",
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
