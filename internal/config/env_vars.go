package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

// LanguageConfig holds the language-specific configuration values
type LanguageConfig struct {
	ImageName        string
	DockerFilePath   string
	TestUserCodePath string
	UtilsFile        string
	AssetsDir        string
}

// Config holds all dynamic configuration values
type ConfigAPI struct {
	MongoDBURI        string
	DBName            string
	Port              string
	FrontendURLS      []string
	Base              string
	LogFilePath       string
	ModeEnv           string
	ClusterName       string
	TemplateAssetsDir string
	Namespace         string
	KubeconfigPath    string
	UniqueAssetsDir   string
	JobTemplatePath   string
}

// NewLanguageConfig creates a new language-specific configuration for a given language.
func newLanguageConfig(language model.PredefinedSupportedLanguage) (*LanguageConfig, error) {
	// Dynamically derive language paths
	// Convert language to lowercase
	TemplateAssetsDir := "./template-assets"
	languageStr := strings.ToLower(string(language))
	langDir := fmt.Sprintf("%s/%s", TemplateAssetsDir, languageStr)
	UtilsFile := fmt.Sprintf("%s/%s/ds_utils.%s", TemplateAssetsDir, strings.ToLower(string(language)), model.GetFileExtension(language))
	dockerFilePath := fmt.Sprintf("%s/Dockerfile", langDir)
	TestUserCodePath := fmt.Sprintf("%s/main.%s", langDir, model.GetFileExtension(language))

	// Validate supported language
	if model.GetFileExtension(language) == "" {
		return nil, model.NewCustomError(400, fmt.Sprintf("unsupported language: %s", language))
	}

	// Create language-specific configuration
	config := &LanguageConfig{
		ImageName:        fmt.Sprintf("skillcode-custom-%s:latest", languageStr),
		DockerFilePath:   dockerFilePath,
		TestUserCodePath: TestUserCodePath,
		UtilsFile:        UtilsFile,
		AssetsDir:        langDir,
	}

	return config, nil
}

// LoadConfigAPI loads the application configuration from environment variables or a config file.
func newConfigAPI() *ConfigAPI {
	return &ConfigAPI{
		MongoDBURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:            getEnv("MONGO_DB", "skillcode_db"),
		Port:              "8080",
		FrontendURLS:      strings.Split(getEnv("FRONTEND_URLS", "http://127.0.0.1:3000,http://127.0.0.1:3001,http://localhost:3000,http://localhost:3001"), ","),
		Base:              "skillcode",
		LogFilePath:       "./logs/app.log",
		ModeEnv:           getEnv("MODE_ENV", "development"),
		ClusterName:       "my-cluster",
		Namespace:         getEnv("NAMESPACE", "default"),
		KubeconfigPath:    getEnv("KUBECONFIG", filepath.Join(os.Getenv("HOME"), ".kube", "config")),
		UniqueAssetsDir:   "./unique-assets",
		JobTemplatePath:   "./template-assets/job-template.yaml",
		TemplateAssetsDir: "./template-assets",
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
	GlobalLanguageConfigs map[model.PredefinedSupportedLanguage]*LanguageConfig
)

func InitGlobalConfigs() error {
	GlobalConfigAPI = newConfigAPI()
	// Initialize GlobalLanguageConfigs map
	GlobalLanguageConfigs = make(map[model.PredefinedSupportedLanguage]*LanguageConfig)

	// Populate the GlobalLanguageConfigs map
	for _, lang := range model.PredefinedSupportedLanguages {
		config, err := newLanguageConfig(lang)
		if err != nil {
			return fmt.Errorf("failed to initialize GlobalLanguageConfig for language %s: %v", lang, err)
		}
		GlobalLanguageConfigs[lang] = config
	}
	return nil
}
