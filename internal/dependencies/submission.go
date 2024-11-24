package dependencies

import (
	"fmt"
	"os"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	tester "github.com/TehilaTheStudent/SkillCode-backend/internal/tester"
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
		return nil, fmt.Errorf("kind not installed: %w", err)
	}

	// 2. Ensure Kind cluster exists
	if err := EnsureClusterExists(config.GlobalConfigAPI.ClusterName); err != nil {
		return nil, fmt.Errorf("failed to create Kind cluster: %w", err)
	}

	// Initialize SharedTester
	sharedTester, err := NewSharedTester(kubeconfigPath, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubernetes client: %w", err)
	}

	return sharedTester, nil
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
	//TODO: its not enough that the cluster exists, it has to be container in docker
	if err != nil {
		return fmt.Errorf("failed to get Kind clusters: %v", err)
	}
	clusters := strings.Split(out, "\n")
	for _, cluster := range clusters {
		if cluster == clusterName {
			return nil // Cluster exists
		}
	}
	_, err = utils.RunCommand("kind", "create", "cluster", "--name", clusterName,"--config", config.GlobalConfigAPI.ClusterConfigFile)
	if err != nil {
		return fmt.Errorf("failed to create Kind cluster: %v", err)
	}
	return nil
}

// NewSharedTester initializes the SharedTester with Kubernetes client
func NewSharedTester(kubeconfigPath, namespace string) (*tester.SharedTester, error) {
	apiServerURL := config.GlobalConfigAPI.KindServerUrl+":"+config.GlobalConfigAPI.ClusterPort
	// Pass the API server URL and kubeconfig path
	config, err := clientcmd.BuildConfigFromFlags(apiServerURL, kubeconfigPath)
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
