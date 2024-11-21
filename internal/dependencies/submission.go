package dependencies

import (
	"fmt"
	"os"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	tester "github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// SetupSubmissionDependencies ensures all dependencies for submission handling are ready
func SetupSubmissionDependencies(kubeconfigPath, namespace string) (*tester.SharedTester, error) {
	// Ensure the working directory is correct
	if err := EnsureWorkingDirectory(); err != nil {
		return nil, fmt.Errorf("working directory setup failed: %w", err)
	}

	// 1. Ensure Kind is installed
	if err := EnsureKindInstalled(); err != nil {
		return nil, fmt.Errorf("failed to set kubectl context: %w", err)
	}

	// 2. Ensure Kind cluster exists
	if err := EnsureClusterExists(config.GlobalConfigAPI.ClusterName); err != nil {
		return nil, fmt.Errorf("failed to create Kind cluster: %w", err)
	}

	// 3. Build Docker images for all languages
	if err := BuildDockerImagesForAllLanguages(); err != nil {
		return nil, fmt.Errorf("failed to build Docker images: %w", err)
	}
	// Initialize SharedTester
	sharedTester, err := NewSharedTester(kubeconfigPath, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubernetes client: %w", err)
	}
	// 4. Load Docker images onto the Kind cluster
	if err := LoadImagesToKindCluster(sharedTester, config.GlobalConfigAPI.ClusterName); err != nil {
		return nil, fmt.Errorf("failed to load images onto the Kind cluster: %w", err)
	}

	return sharedTester, nil
}

// BuildDockerImagesForAllLanguages builds Docker images for all predefined supported languages
func BuildDockerImagesForAllLanguages() error {
	for _, lang := range model.PredefinedSupportedLanguages {
		langConfig, exists := config.GlobalLanguageConfigs[lang]
		if !exists {
			return fmt.Errorf("configuration for language %s not found", lang)
		}
		fmt.Printf("Building Docker image '%s' with Dockerfile '%s' in context '%s'...\n", langConfig.ImageName, langConfig.DockerFilePath, config.GlobalConfigAPI.TemplateAssetsDir)
		output, err := utils.RunCommand(
			"docker",
			"build",
			"-t", langConfig.ImageName, // Tag the image
			"-f", langConfig.DockerFilePath, // Specify the Dockerfile
			config.GlobalConfigAPI.TemplateAssetsDir, // Use the language-specific directory as the build context
		)

		if err != nil {
			return fmt.Errorf("failed to build Docker image '%s': %v\nOutput: %s", langConfig.ImageName, err, output)
		}
	}
	return nil
}

// LoadImagesToKindCluster loads all built images into the Kind cluster using kind commands
func LoadImagesToKindCluster(t *tester.SharedTester, clusterName string) error {
	// Loop through all predefined supported languages
	for _, lang := range model.PredefinedSupportedLanguages {
		config, exists := config.GlobalLanguageConfigs[lang]
		if !exists {
			return fmt.Errorf("configuration for language %s not found", lang)
		}

		// Check if the image exists locally
		imageName := config.ImageName


		// Load the image into the Kind cluster using kind command
		fmt.Printf("Loading image '%s' into Kind cluster '%s'...\n", imageName, clusterName)
		_, err := utils.RunCommand("kind", "load", "docker-image", imageName, "--name", clusterName)
		if err != nil {
			return fmt.Errorf("failed to load image '%s' into Kind cluster: %v", imageName, err)
		}

		fmt.Printf("Successfully loaded image '%s' into Kind cluster\n", imageName)
	}
	return nil
}

func EnsureKindInstalled() error {
	out, err := utils.RunCommand("kind", "--version")
	if err != nil {
		return fmt.Errorf("Kind is not installed or misconfigured: %v", err)
	}
	fmt.Printf("Kind version: %s\n", out)
	return nil
}

func EnsureClusterExists(clusterName string) error {
	out, err := utils.RunCommand("kind", "get", "clusters")
	if err != nil {
		return fmt.Errorf("failed to get Kind clusters: %v", err)
	}
	clusters := strings.Split(out, "\n")
	for _, cluster := range clusters {
		if cluster == clusterName {
			return nil // Cluster exists
		}
	}
	_, err = utils.RunCommand("kind", "create", "cluster", "--name", clusterName)
	if err != nil {
		return fmt.Errorf("failed to create Kind cluster: %v", err)
	}
	return nil
}

// NewSharedTester initializes the SharedTester with Kubernetes client
func NewSharedTester(kubeconfigPath, namespace string) (*tester.SharedTester, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Create SharedTester
	return &tester.SharedTester{
		ClientSet: clientset,
		Namespace: namespace,
	}, nil
}

// EnsureWorkingDirectory ensures the working directory is set to the project root specified by an environment variable
func EnsureWorkingDirectory() error {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		return fmt.Errorf("PROJECT_ROOT environment variable is not set")
	}
	if err := os.Chdir(projectRoot); err != nil {
		return fmt.Errorf("failed to change directory to PROJECT_ROOT: %w", err)
	}
	return nil
}

func Cleanup() error {
	// 1. Remove temporary directories
	tempDir := config.GlobalConfigAPI.UniqueAssetsDir
	err := os.RemoveAll(tempDir)
	if err != nil {
		return fmt.Errorf("failed to remove temporary directories: %v", err)
	}

	// 2. Remove Docker containers/images
	for _, config := range config.GlobalLanguageConfigs {
		_, err = utils.RunCommand("docker", "rmi", "-f", config.ImageName)
		if err != nil {
			return fmt.Errorf("failed to remove Docker image '%s': %v", config.ImageName, err)
		}
	}

	// 3. Remove Kind cluster
	clusterName := config.GlobalConfigAPI.ClusterName
	_, err = utils.RunCommand("kind", "delete", "cluster", "--name", clusterName)
	if err != nil {
		return fmt.Errorf("failed to delete Kind cluster '%s': %v", clusterName, err)
	}

	return nil
}
