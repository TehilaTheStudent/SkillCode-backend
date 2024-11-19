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
	ImageName        string
	DockerFilePath   string
	PodFilePath      string
	PodName          string
	ClusterName      string
	TestUserCodePath string
	UtilsFile        string
	BaseDir          string
}

// NewSandboxConfig creates a new sandbox configuration for a given language.
func newSandboxConfig(language model.PredefinedSupportedLanguage) (*ConfigSandbox, error) {
	// Dynamically derive language paths
	// Convert language to lowercase
	BaseDir := "./assets"
	languageStr := strings.ToLower(string(language))
	langDir := fmt.Sprintf("%s/%s", BaseDir, languageStr)
	defaultPodFilePath := fmt.Sprintf("%s/pod.yaml", langDir)
	UtilsFile := BaseDir + strings.ToLower(string(language)) + "/ds_utils." + model.GetFileExtension(language)
	dockerFilePath := fmt.Sprintf("%s/Dockerfile", langDir)
	TestUserCodePath := fmt.Sprintf("%s/run_tests.%s", langDir, model.GetFileExtension(language))

	// Validate supported language
	if model.GetFileExtension(language) == "" {
		return nil, model.NewCustomError(400, fmt.Sprintf("unsupported language: %s", language))
	}

	// Create sandbox configuration
	config := &ConfigSandbox{
		ImageName:        fmt.Sprintf("%s-test:latest", languageStr),
		DockerFilePath:   dockerFilePath,
		PodFilePath:      defaultPodFilePath,
		TestUserCodePath: TestUserCodePath,
		UtilsFile:        UtilsFile,
		BaseDir:          BaseDir,
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
	LogFilePath string
	ModeEnv     string
	PodName     string
	ClusterName string
}

// LoadConfigAPI loads the application configuration from environment variables or a config file.
func newConfigAPI() *ConfigAPI {
	return &ConfigAPI{
		MongoDBURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:      getEnv("MONGO_DB", "skillcode_db"),
		Port:        getEnv("PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://127.0.0.1:3000"),
		Base:        "skillcode",
		LogFilePath: "./logs/app.log",
		ModeEnv:     getEnv("MODE_ENV", "development"),
		PodName:     "pod",
		ClusterName: "my-cluster",
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

var (
	GlobalConfigAPI       *ConfigAPI
	GlobalConfigSandboxes map[model.PredefinedSupportedLanguage]*ConfigSandbox
)

func InitGlobalConfigs() error {
	GlobalConfigAPI = newConfigAPI()
	// Initialize GlobalConfigSandbox map
	GlobalConfigSandboxes = make(map[model.PredefinedSupportedLanguage]*ConfigSandbox)

	// Populate the GlobalConfigSandbox map
	for _, lang := range model.PredefinedSupportedLanguages {
		config, err := newSandboxConfig(lang)
		if err != nil {
			return fmt.Errorf("failed to initialize GlobalConfigSandbox for language %s: %v", lang, err)
		}
		GlobalConfigSandboxes[lang] = config
	}
	return nil
}
