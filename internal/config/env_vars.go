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

type ConfigCode struct {
	UtilsFile string
}

const BaseDir = "./assets/"

func NewConfigCode(language model.PredefinedSupportedLanguage) *ConfigCode {
	return &ConfigCode{
		UtilsFile: BaseDir + strings.ToLower(string(language)) + "/ds_utils",
	}
}

// NewSandboxConfig creates a new sandbox configuration for a given language.
func NewSandboxConfig(language model.PredefinedSupportedLanguage) (*ConfigSandbox, error) {
	// Dynamically derive language paths
	// Convert language to lowercase
	languageStr := strings.ToLower(string(language))
	langDir := fmt.Sprintf("%s/%s", BaseDir, languageStr)
	defaultPodFilePath := fmt.Sprintf("%s/pod.yaml", langDir)
	clusterName := "my-cluster"

	dockerFilePath := fmt.Sprintf("%s/Dockerfile", langDir)
	userCodePath := fmt.Sprintf("%s/run_tests.%s", langDir, model.GetFileExtension(language))

	// Validate supported language
	if model.GetFileExtension(language) == "" {
		return nil, model.NewCustomError(400, fmt.Sprintf("unsupported language: %s", language))
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
		FrontendURL: getEnv("FRONTEND_URL", "http://127.0.0.1:3000"),
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
