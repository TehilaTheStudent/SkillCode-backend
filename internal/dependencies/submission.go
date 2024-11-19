package dependencies

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
	"go.uber.org/zap"
)

// EnsureWorkingDirectory ensures the working directory is set to the project root specified by an environment variable
func EnsureWorkingDirectory() {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		log.Println("PROJECT_ROOT environment variable is not set")
		os.Exit(1)
	}
	if err := changeDirectory(projectRoot); err != nil {
		log.Fatalf("Failed to change directory to PROJECT_ROOT: %v", err)
	}
}

// changeDirectory changes the current working directory to the specified path
func changeDirectory(path string) error {
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", path, err)
	}
	return nil
}

func SetupDependencies(config config.ConfigSandbox, logger *zap.Logger) error {
	logger.Info("Starting pre-startup setup")

	// 1. Check if kubectl is installed
	logger.Info("Checking if kubectl is installed")
	if err := EnsureKubectlInstalled(); err != nil {
		logger.Error("Kubectl check failed", zap.Error(err))
		return fmt.Errorf("kubectl check failed: %w", err)
	}

	// 2. Ensure Kind cluster exists
	logger.Info("Ensuring Kind cluster exists", zap.String("cluster", config.ClusterName))
	if err := EnsureClusterExists(config.ClusterName); err != nil {
		logger.Error("Cluster setup failed", zap.String("cluster", config.ClusterName), zap.Error(err))
		return fmt.Errorf("failed to create Kind cluster: %w", err)
	}

	// 3. Set kubectl context
	logger.Info("Setting kubectl context", zap.String("cluster", config.ClusterName))
	if err := EnsureKindContext(config.ClusterName); err != nil {
		logger.Error("Failed to set kubectl context", zap.Error(err))
		return fmt.Errorf("failed to set kubectl context: %w", err)
	}

	logger.Info("Pre-startup setup completed successfully")
	return nil
}

func EnsureKindContext(clusterName string) error {
	_, err :=  utils.RunCommand("kubectl", "config", "use-context", fmt.Sprintf("kind-%s", clusterName))
	if err != nil {
		return model.NewCustomError(500, fmt.Sprintf("failed to switch kubectl context: %v", err))
	}
	return nil
}

func DeleteKindCluster(clusterName string) error {
	_, err :=  utils.RunCommand("kind", "delete", "cluster", "--name", clusterName)
	if err != nil {
		return model.NewCustomError(500, fmt.Sprintf("failed to delete Kind cluster: %v", err))
	}
	return nil
}

func EnsureKubectlInstalled() error {
	_, err :=  utils.RunCommand("kubectl", "version", "--client")
	if err != nil {
		return model.NewCustomError(500, fmt.Sprintf("kubectl not installed or misconfigured: %v", err))
	}
	return nil
}


func EnsureClusterExists(clusterName string) error {
	out, err :=  utils.RunCommand("kind", "get", "clusters")
	if err != nil {
		return model.NewCustomError(500, fmt.Sprintf("failed to get Kind clusters: %v", err))
	}
	clusters := strings.Split(out, "\n")
	for _, cluster := range clusters {
		if cluster == clusterName {
			return nil // Cluster exists
		}
	}
	_, err =  utils.RunCommand("kind", "create", "cluster", "--name", clusterName)
	if err != nil {
		return model.NewCustomError(500, fmt.Sprintf("failed to create Kind cluster: %v", err))
	}
	return nil
}
