package utils

import (
	"fmt"
	"log"
	"os"
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
