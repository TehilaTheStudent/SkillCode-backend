package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func InitLogger() *zap.Logger {
	// Ensure the logs directory exists
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create logs directory: %v", err))
	}

	// Create or open the log file
	logFile := fmt.Sprintf("%s/app.log", logsDir)
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}

	// Configure file writer and console writer
	fileWriteSyncer := zapcore.AddSync(file)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	// Create an encoder (format for log entries)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create core with file and console writers
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                            // JSON formatted logs
		zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer), // Log to both file and console
		zapcore.InfoLevel, // Log level
	)

	// Create the logger
	logger := zap.New(core, zap.AddCaller())

	return logger
}
